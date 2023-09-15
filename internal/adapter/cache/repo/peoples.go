package repo

import (
	"context"
	"fmt"
	"github.com/Dmitrij-Kochetov/peoples/internal/domain/dto"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"time"
)

type CachePeopleRepo struct {
	Client *redis.Client
	exp    time.Duration
}

func NewCachePeopleRepo(client *redis.Client, exp time.Duration) *CachePeopleRepo {
	return &CachePeopleRepo{Client: client, exp: exp}
}

func (c *CachePeopleRepo) Create(ctx context.Context, people dto.People) error {
	data, err := people.MarshallBinary()
	if err != nil {
		return err
	}

	if err := c.Client.Set(ctx, people.ID.String(), data, c.exp).Err(); err != nil {
		return err
	}

	return nil
}

func (c *CachePeopleRepo) FindById(ctx context.Context, uuid uuid.UUID) (*dto.People, error) {
	result, err := c.Client.Get(ctx, uuid.String()).Result()
	if err != nil {
		return nil, err
	}

	if result == "" {
		return nil, fmt.Errorf("NotFound")
	}

	var people dto.People
	if err := people.UnmarshallBinary([]byte(result)); err != nil {
		return nil, err
	}
	return &people, nil
}

func (c *CachePeopleRepo) Update(ctx context.Context, people dto.People) error {
	if err := c.Delete(ctx, people.ID); err != nil {
		return err
	}

	return c.Create(ctx, people)
}

func (c *CachePeopleRepo) Delete(ctx context.Context, uuid uuid.UUID) error {
	if err := c.Client.Del(ctx, uuid.String()).Err(); err != nil {
		return err
	}
	return nil
}
