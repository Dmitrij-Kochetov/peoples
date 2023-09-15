package repo

import (
	"context"
	cache "github.com/Dmitrij-Kochetov/peoples/internal/adapter/cache/repo"
	db "github.com/Dmitrij-Kochetov/peoples/internal/adapter/database/repo"
	"github.com/Dmitrij-Kochetov/peoples/internal/domain/dto"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"time"
)

type PeopleRepo struct {
	db    *db.DbPeopleRepo
	cache *cache.CachePeopleRepo
}

func NewPeopleRepo(dbConn *sqlx.DB, client *redis.Client, exp time.Duration) *PeopleRepo {
	return &PeopleRepo{
		db:    db.NewDbPeopleRepo(dbConn),
		cache: cache.NewCachePeopleRepo(client, exp),
	}
}

func (p *PeopleRepo) GetByID(ctx context.Context, uuid uuid.UUID) (*dto.People, error) {
	res, err := p.cache.FindById(ctx, uuid)
	if err == redis.Nil {
		res, err := p.db.GetByID(uuid)
		if err != nil {
			return nil, err
		}
		if err := p.cache.Create(ctx, *res); err != nil {
			return nil, err
		}
		return res, nil
	}
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (p *PeopleRepo) GetAllByFilter(ctx context.Context, filter dto.Filter) (*dto.Peoples, error) {
	return p.db.GetAllByFilter(filter)
}

func (p *PeopleRepo) Create(ctx context.Context, people dto.CreatePeople) error {
	return p.db.Create(people)
}

func (p *PeopleRepo) Update(ctx context.Context, people dto.People) error {
	if err := p.db.Update(people); err != nil {
		return err
	}

	return p.cache.Update(ctx, people)
}

func (p *PeopleRepo) DeleteByID(ctx context.Context, uuid uuid.UUID) error {
	if err := p.db.DeleteByID(uuid); err != nil {
		return nil
	}

	return p.cache.Delete(ctx, uuid)
}

func (p *PeopleRepo) Close(ctx context.Context) error {
	if err := p.db.DB.Close(); err != nil {
		return err
	}
	if err := p.cache.Client.Close(); err != nil {
		return err
	}
	return nil
}
