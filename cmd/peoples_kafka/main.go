package main

import (
	"context"
	"github.com/Dmitrij-Kochetov/peoples/internal/adapter/config"
	"github.com/Dmitrij-Kochetov/peoples/internal/adapter/config/kafka_config"
	"github.com/Dmitrij-Kochetov/peoples/internal/application/presentation/kafka"
	_ "github.com/lib/pq"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := config.LoadConfig(kafka_config.Config{})
	errChan, err := run(cfg)
	if err != nil {
		log.Fatalf("Couldn't run: %v", err)
	}
	if err := <-errChan; err != nil {
		log.Fatalf("Error while running: %v", err)
	}
}

func run(cfg kafka_config.Config) (<-chan error, error) {
	server, err := kafka.NewServerFromConfig(cfg)
	if err != nil {
		return nil, err
	}

	errChan := make(chan error, 1)

	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	go func() {
		<-ctx.Done()
		log.Println("Shutting down...")

		ctxTimeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		defer func() {
			err := server.Close()
			if err != nil {
				log.Fatalf("Couldn't close server: %v", err)
			}
			stop()
			cancel()
			close(errChan)
		}()

		if err := server.Shutdown(ctxTimeout); err != nil {
			errChan <- err
		}

		log.Println("Gracefully shutting down")
	}()

	go func() {
		log.Println("Starting server...")

		if err := server.ListenAndServe(); err != nil {
			errChan <- err
		}
	}()

	return errChan, nil
}
