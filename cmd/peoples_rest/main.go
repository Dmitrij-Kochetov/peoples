package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Dmitrij-Kochetov/peoples/internal/adapter/config"
	"github.com/Dmitrij-Kochetov/peoples/internal/adapter/config/rest_config"
	"github.com/Dmitrij-Kochetov/peoples/internal/application/presentation/rest"
)

func main() {
	var cfg rest_config.Config
	cfg = config.LoadConfig(cfg)
	errChan, err := run(cfg)
	if err != nil {
		log.Fatalf("Couldn't run: %v", err)
	}
	if err := <-errChan; err != nil {
		log.Fatalf("Error while running: %v", err)
	}
}

func run(cfg rest_config.Config) (<-chan error, error) {
	server, err := rest.NewServerFromConfig(cfg)
	if err != nil {
		return nil, err
	}

	errChan := make(chan error, 1)

	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	httpServ := server.GetHttp()

	go func() {
		<-ctx.Done()
		log.Println("Shutting down...")

		ctxTimeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		defer func() {
			redisCtxTimeout, c := context.WithTimeout(ctxTimeout, 5*time.Second)
			err := server.Close(redisCtxTimeout)
			if err != nil {
				log.Fatalf("Couldn't close server: %v", err)
			}
			stop()
			cancel()
			c()
			close(errChan)
		}()

		if err := httpServ.Shutdown(ctxTimeout); err != nil {
			errChan <- err
		}

		if err := server.Shutdown(ctxTimeout); err != nil {
			errChan <- err
		}

		log.Println("Gracefully shutting down")
	}()

	go func() {
		log.Println("Starting server...")

		if err := httpServ.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}

	}()

	return errChan, nil
}
