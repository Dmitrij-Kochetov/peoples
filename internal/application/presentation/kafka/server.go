package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"

	"github.com/Dmitrij-Kochetov/peoples/internal/adapter/config/kafka_config"
	db "github.com/Dmitrij-Kochetov/peoples/internal/adapter/database/repo"
	internal "github.com/Dmitrij-Kochetov/peoples/internal/adapter/kafka"
	"github.com/Dmitrij-Kochetov/peoples/internal/adapter/logging"
	"github.com/Dmitrij-Kochetov/peoples/internal/application/usecases"
	dto "github.com/Dmitrij-Kochetov/peoples/internal/domain/dto/kafka"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/jmoiron/sqlx"
)

type Server struct {
	logger     *slog.Logger
	consumer   *internal.Consumer
	producer   *internal.Producer
	peopleRepo *db.DbPeopleRepo
	doneChan   chan struct{}
	closeChan  chan struct{}
}

func NewServerFromConfig(config kafka_config.Config) (*Server, error) {
	logger := logging.SetUpLogger(config.Env)

	dbConn, err := sqlx.Connect(config.DB.Driver, config.DB.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect %w", err)
	}

	peopleRepo := db.NewDbPeopleRepo(dbConn)

	consumer, err := internal.NewKafkaConsumer(config.Kafka.Address,
		config.Kafka.ConsumerTopic,
		config.Kafka.ConsumerGroup,
		config.Kafka.Timeout,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka consumer %w", err)
	}

	producer, err := internal.NewKafkaProducer(config.Kafka.Address, config.Kafka.ProducerTopic)

	if err != nil {
		return nil, fmt.Errorf("failed to create kafka producer %w", err)
	}

	return &Server{
		logger:     logger,
		consumer:   consumer,
		producer:   producer,
		peopleRepo: peopleRepo,
		doneChan:   make(chan struct{}),
		closeChan:  make(chan struct{}),
	}, nil
}

func (s *Server) ListenAndServe() error {
	commit := func(msg *kafka.Message) {
		if _, err := s.consumer.Consumer.CommitMessage(msg); err != nil {
			s.logger.Error("commit failed", logging.Err(err))
		}
	}

	writeError := func(e dto.Error) error {
		payloadBytes, err := json.Marshal(e)
		if err != nil {
			return err
		}

		err = s.producer.Producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{
				Topic:     &s.producer.Topic,
				Partition: kafka.PartitionAny,
			},
			Value: payloadBytes,
		}, nil)
		if err != nil {
			return err
		}
		return nil
	}

	handleError := func(e dto.Error) {
		s.logger.Error(e.Message, slog.Attr{
			Key:   "error",
			Value: slog.StringValue(e.Error),
		})

		err := writeError(e)
		if err != nil {
			s.logger.Error("failed to write error to kafka", logging.Err(err))
		}
	}

	go func() {
		run := true

		wg := sync.WaitGroup{}

		for run {
			select {
			case <-s.closeChan:
				run = false
				wg.Wait()
				break
			default:
				msg, ok := s.consumer.Consumer.Poll(s.consumer.TimeoutMs).(*kafka.Message)
				if !ok {
					continue
				}
				s.logger.Info("message received")

				var payload dto.PeopleName
				if err := json.Unmarshal(msg.Value, &payload); err != nil {
					go func() {
						wg.Add(1)
						defer wg.Done()
						handleError(dto.Error{
							Message: "failed to unmarshal payload",
							Error:   err.Error(),
						})
					}()
					commit(msg)
					continue
				}

				if payload.FirstName == nil || payload.LastName == nil {
					go func() {
						wg.Add(1)
						defer wg.Done()
						handleError(dto.ErrRequiredFieldNotExists)
					}()

					commit(msg)
					continue
				}

				if *payload.FirstName == "" || *payload.LastName == "" {
					go func() {
						wg.Add(1)
						defer wg.Done()
						handleError(dto.ErrRequiredFieldIsEmpty)
					}()

					commit(msg)
					continue
				}

				go func() {
					wg.Add(1)
					defer wg.Done()
					agifyInfo, err := usecases.AgifyPeople(*payload.FirstName)
					if err != nil {
						handleError(dto.Error{
							Message: err.Error(),
							Error:   "agified failed",
						})
						commit(msg)
						return
					}
					err = usecases.CreateAgifiedPeople(s.peopleRepo, payload, agifyInfo)
					if err != nil {
						handleError(dto.Error{
							Message: err.Error(),
							Error:   "create agified failed",
						})
						commit(msg)
						return
					}
					commit(msg)
					return
				}()
			}
		}
		s.logger.Info("consumer stopped")
		s.doneChan <- struct{}{}
	}()

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("shutting down")
	close(s.closeChan)

	for {
		select {
		case <-s.doneChan:
			return nil
		case <-ctx.Done():
			return fmt.Errorf("context canceled: %w", ctx.Err())
		}
	}
}

func (s *Server) Close() error {
	if err := s.consumer.Consumer.Unsubscribe(); err != nil {
		return err
	}
	s.producer.Producer.Close()
	if err := s.peopleRepo.DB.Close(); err != nil {
		return err
	}
	return nil
}
