package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Dmitrij-Kochetov/peoples/internal/adapter/config/kafka_config"
	db "github.com/Dmitrij-Kochetov/peoples/internal/adapter/database/repo"
	"github.com/Dmitrij-Kochetov/peoples/internal/adapter/logging"
	"github.com/Dmitrij-Kochetov/peoples/internal/application/usecases"
	dto "github.com/Dmitrij-Kochetov/peoples/internal/domain/dto/kafka"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/jmoiron/sqlx"
	"log/slog"
	"sync"
)

type Server struct {
	logger     *slog.Logger
	consumer   *Consumer
	producer   *Producer
	peopleRepo *db.DbPeopleRepo
	doneChan   chan struct{}
	closeChan  chan struct{}
}

func NewServerFromConfig(config kafka_config.Config) (*Server, error) {
	logger := logging.SetUpLogger(config.Env)

	dbConn, err := sqlx.Connect(config.DB.Driver, config.DB.DbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect %w", err)
	}

	peopleRepo := db.NewDbPeopleRepo(dbConn)

	consumer, err := NewKafkaConsumer(config.Kafka.KafkaURL,
		config.Kafka.ConsumerTopic,
		config.Kafka.ConsumerGroup,
		config.Kafka.KafkaTimeout,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka consumer %w", err)
	}

	producer, err := NewKafkaProducer(config.Kafka.KafkaURL, config.Kafka.ProducerTopic)

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
					s.logger.Error("failed to unmarshal payload", logging.Err(err))

					go func() {
						wg.Add(1)
						defer wg.Done()
						if err := writeError(dto.Error{
							Message: "failed to unmarshal payload",
							Error:   err.Error(),
						}); err != nil {
							s.logger.Error("failed to write error", logging.Err(err))
						}
					}()

					commit(msg)
					continue
				}

				if payload.FirstName == nil || payload.LastName == nil {
					s.logger.Error("invalid payload: first name or last name is nil")

					go func() {
						wg.Add(1)
						defer wg.Done()
						if err := writeError(dto.Error{
							Message: "failed to unmarshal payload",
							Error:   fmt.Errorf("invalid payload: first name or last name is nil").Error(),
						}); err != nil {
							s.logger.Error("failed to write error", logging.Err(err))
						}
					}()

					commit(msg)
					continue
				}

				go func() {
					wg.Add(1)
					defer wg.Done()
					agifyInfo, err := usecases.AgifyPeople(*payload.FirstName)
					if err != nil {
						s.logger.Error("failed to get agify info", logging.Err(err))
						return
					}
					err = usecases.CreateAgifiedPeople(s.peopleRepo, payload, agifyInfo)
					if err != nil {
						s.logger.Error("failed to create agified people", logging.Err(err))
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

func (s *Server) CloseAll() {
	s.consumer.Consumer.Unsubscribe()
	s.producer.Producer.Close()
	s.peopleRepo.DB.Close()
}
