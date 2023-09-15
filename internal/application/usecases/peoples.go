package usecases

import (
	"context"
	"github.com/Dmitrij-Kochetov/peoples/internal/domain/dto"
	"github.com/google/uuid"
)

type IPeopleRepo interface {
	GetByID(context.Context, uuid.UUID) (*dto.People, error)
	GetAllByFilter(context.Context, dto.Filter) (*dto.Peoples, error)
	Create(context.Context, dto.CreatePeople) error
	Update(context.Context, dto.People) error
	DeleteByID(context.Context, uuid.UUID) error
}

func GetPeopleByID(ctx context.Context, repo IPeopleRepo, id uuid.UUID) (*dto.People, error) {
	return repo.GetByID(ctx, id)
}

func GetAllPeopleByFilter(ctx context.Context, repo IPeopleRepo, filter dto.Filter) (*dto.Peoples, error) {
	return repo.GetAllByFilter(ctx, filter)
}

func CreatePeople(ctx context.Context, repo IPeopleRepo, people dto.CreatePeople) error {
	return repo.Create(ctx, people)
}

func UpdatePeopleByID(ctx context.Context, repo IPeopleRepo, people dto.People) error {
	return repo.Update(ctx, people)
}

func DeletePeopleByID(ctx context.Context, repo IPeopleRepo, id uuid.UUID) error {
	return repo.DeleteByID(ctx, id)
}
