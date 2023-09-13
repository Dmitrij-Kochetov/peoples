package repo

import (
	"github.com/Dmitrij-Kochetov/peoples/internal/domain/dto"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type DbPeopleRepo struct {
	db *sqlx.DB
}

func NewDbPeopleRepo(conn *sqlx.DB) *DbPeopleRepo {
	return &DbPeopleRepo{db: conn}
}

func (p *DbPeopleRepo) GetByID(uuid uuid.UUID) (*dto.People, error) {
	var people dto.People

	if err := p.db.Get(&people,
		`SELECT * FROM peoples WHERE id=$1`,
		uuid,
	); err != nil {
		return nil, err
	}

	return &people, nil
}

func (p *DbPeopleRepo) GetAllByFilter(filter dto.Filter) (*dto.Peoples, error) {
	var peoples dto.Peoples

	if err := p.db.Select(&peoples,
		`SELECT * FROM peoples WHERE deleted=$1 ORDER BY id LIMIT $2 OFFSET $3`,
		filter.Deleted,
		filter.Limit,
		filter.Offset,
	); err != nil {
		return nil, err
	}

	return &peoples, nil
}

func (p *DbPeopleRepo) Create(people dto.CreatePeople) (*uuid.UUID, error) {
	tx, err := p.db.Begin()
	if err != nil {
		return nil, err
	}

	var uuidString string

	if _, err := tx.Exec(
		`INSERT INTO peoples (first_name, last_name, patronymic, age, sex, nation)
			VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		people.FirstName,
		people.LastName,
		people.Patronymic,
		people.Age,
		people.Sex,
		people.Nation,
		&uuidString,
	); err != nil {
		if err := tx.Rollback(); err != nil {
			return nil, err
		}
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	parsed, err := uuid.Parse(uuidString)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}

func (p *DbPeopleRepo) Update(people dto.People) error {
	tx, err := p.db.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.Exec(
		`UPDATE peoples 
			SET first_name=$1, second_name=$2, patronymic=$3, age=$4, sex=$5, nation=$6, deleted=$7
			WHERE id=$8`,
		people.FirstName,
		people.LastName,
		people.Patronymic,
		people.Age,
		people.Sex,
		people.Nation,
		people.Deleted,
		people.ID,
	); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (p *DbPeopleRepo) DeleteByID(uuid uuid.UUID) error {
	tx, err := p.db.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.Exec(
		`UPDATE peoples SET deleted=$1 WHERE id=$2`,
		true,
		uuid,
	); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}