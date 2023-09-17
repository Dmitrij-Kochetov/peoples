package rest

import (
	"net/http"

	"github.com/Dmitrij-Kochetov/peoples/internal/domain/dto"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type FilterRequest struct {
	Limit   int  `json:"limit"`
	Offset  int  `json:"offset"`
	Deleted bool `json:"deleted"`
}

func (f *FilterRequest) Bind(r *http.Request) error {
	return nil
}

type PeopleResponse struct {
	ID         uuid.UUID `json:"id" db:"id"`
	FirstName  string    `json:"first_name" db:"first_name"`
	LastName   string    `json:"last_name" db:"last_name"`
	Patronymic string    `json:"patronymic" db:"patronymic"`
	Age        int       `json:"age" db:"age"`
	Sex        string    `json:"sex" db:"sex"`
	Nation     string    `json:"nation" db:"nation"`
	Deleted    bool      `json:"deleted" db:"deleted"`
}

type CreatePeopleRequest struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Patronymic string `json:"patronymic"`
	Age        int    `json:"age"`
	Sex        string `json:"sex"`
	Nation     string `json:"nation"`
}

func (c *CreatePeopleRequest) Bind(r *http.Request) error {
	return nil
}

type UpdatePeopleRequest struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Patronymic string `json:"patronymic"`
	Age        int    `json:"age"`
	Sex        string `json:"sex"`
	Nation     string `json:"nation"`
}

func (u *UpdatePeopleRequest) Bind(r *http.Request) error {
	return nil
}

func NewPeopleResponse(people dto.People) *PeopleResponse {
	return &PeopleResponse{
		ID:         people.ID,
		FirstName:  people.FirstName,
		LastName:   people.LastName,
		Patronymic: people.Patronymic,
		Age:        people.Age,
		Sex:        people.Sex,
		Deleted:    people.Deleted,
	}
}

func (*PeopleResponse) Render(w http.ResponseWriter, req *http.Request) error {
	return nil
}

func NewListPeopleResponse(peoples *dto.Peoples) []render.Renderer {
	r := make([]render.Renderer, len(*peoples))
	for idx, people := range *peoples {
		r[idx] = NewPeopleResponse(people)
	}
	return r
}
