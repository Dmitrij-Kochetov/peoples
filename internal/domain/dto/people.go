package dto

import (
	"encoding/json"

	"github.com/google/uuid"
)

type People struct {
	ID         uuid.UUID `json:"id" db:"id"`
	FirstName  string    `json:"first_name" db:"first_name"`
	LastName   string    `json:"last_name" db:"last_name"`
	Patronymic string    `json:"patronymic" db:"patronymic"`
	Age        int       `json:"age" db:"age"`
	Sex        string    `json:"sex" db:"sex"`
	Nation     string    `json:"nation" db:"nation"`
	Deleted    bool      `json:"deleted" db:"deleted"`
}

type CreatePeople struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Patronymic string `json:"patronymic,omitempty"`
	Age        int    `json:"age,omitempty"`
	Sex        string `json:"sex,omitempty"`
	Nation     string `json:"nation,omitempty"`
}

type Peoples []People

func (r *People) MarshallBinary() ([]byte, error) {
	return json.Marshal(r)
}

func (r *People) UnmarshallBinary(data []byte) error {
	if err := json.Unmarshal(data, &r); err != nil {
		return err
	}
	return nil
}
