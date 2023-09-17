package rest

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/Dmitrij-Kochetov/peoples/internal/adapter/config/rest_config"
	"github.com/Dmitrij-Kochetov/peoples/internal/adapter/logging"
	"github.com/Dmitrij-Kochetov/peoples/internal/adapter/repo"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

type serverCfg struct {
	addr        string
	timeout     time.Duration
	idleTimeout time.Duration
}

type Server struct {
	logger    *slog.Logger
	repo      *repo.PeopleRepo
	router    *chi.Mux
	cfg       serverCfg
	doneChan  chan struct{}
	closeChan chan struct{}
}

func NewServerFromConfig(cfg rest_config.Config) (*Server, error) {
	logger := logging.SetUpLogger(cfg.Env)

	dbConn, err := sqlx.Connect(cfg.Db.Driver, cfg.Db.Url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect %w", err)
	}

	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Address,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping redis %w", err)
	}

	repos := repo.NewPeopleRepo(dbConn, client, cfg.Redis.Timeout)

	return &Server{
		logger:    logger,
		repo:      repos,
		router:    chi.NewRouter(),
		doneChan:  make(chan struct{}),
		closeChan: make(chan struct{}),
		cfg: serverCfg{
			addr:        cfg.Server.Address,
			timeout:     cfg.Server.Timeout,
			idleTimeout: cfg.Server.IdleTimeout,
		},
	}, nil
}

func (s *Server) GetHttp() *http.Server {
	s.setupRoutes()
	srv := http.Server{
		Addr:         s.cfg.addr,
		Handler:      s.router,
		IdleTimeout:  s.cfg.idleTimeout,
		ReadTimeout:  s.cfg.timeout,
		WriteTimeout: s.cfg.timeout,
	}

	return &srv
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("shutting down")
	close(s.closeChan)

	for {
		select {
		case <-s.closeChan:
			return nil
		case <-ctx.Done():
			return fmt.Errorf("context canceled: %w", ctx.Err())
		}
	}
}

func (s *Server) Close(ctx context.Context) error {
	if err := s.repo.Close(ctx); err != nil {
		return err
	}
	return nil
}
