package repo

import (
	"context"
	"github.com/Dmitrij-Kochetov/peoples/internal/domain/dto"
	"github.com/google/uuid"
)

type IPeopleRepo interface {
	GetByID(context.Context, uuid.UUID) (dto.People, error)
	GetAllByFilter(context.Context, dto.Filter) (dto.Peoples, error)
	Create(context.Context, dto.CreatePeople) (uuid.UUID, error)
	Update(context.Context, dto.People) error
	DeleteByID(context.Context, uuid.UUID) error
}
