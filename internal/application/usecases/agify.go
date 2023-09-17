package usecases

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	db "github.com/Dmitrij-Kochetov/peoples/internal/adapter/database/repo"
	"github.com/Dmitrij-Kochetov/peoples/internal/domain/dto"
	"github.com/Dmitrij-Kochetov/peoples/internal/domain/dto/kafka"
)

type AgifyResponse struct {
	Age int `json:"age"`
}

type GenderizeResponse struct {
	Gender string `json:"gender"`
}

type NationalizeResponse struct {
	Country []nation `json:"country"`
}

type nation struct {
	CountryId string `json:"country_id"`
}

type Responses interface {
	AgifyResponse | GenderizeResponse | NationalizeResponse
}

type AgifyInfo struct {
	Age    int
	Sex    string
	Nation string
}

func doRequest[R Responses](url, name string, response R) (R, error) {
	resp, err := http.Get(url + name)
	if err != nil {
		return *new(R), err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return *new(R), err
	}

	if err := resp.Body.Close(); err != nil {
		return *new(R), err
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return *new(R), err
	}
	return response, nil
}

func AgifyPeople(name string) (AgifyInfo, error) {
	var urls = []string{"https://api.agify.io/?name=", "https://api.genderize.io/?name=", "https://api.nationalize.io/?name="}

	age, err := doRequest(urls[0], name, AgifyResponse{})
	if err != nil {
		return AgifyInfo{}, err
	}
	sex, err := doRequest(urls[1], name, GenderizeResponse{})
	if err != nil {
		return AgifyInfo{}, err
	}
	if sex.Gender == "" {
		return AgifyInfo{}, fmt.Errorf("cannot get existing gender, possibly name is wrong")
	}

	nation, err := doRequest(urls[2], name, NationalizeResponse{})
	if err != nil {
		return AgifyInfo{}, err
	}

	return AgifyInfo{
		Age: age.Age,
		Sex: sex.Gender,
		Nation: func() string {
			if len(nation.Country) == 0 {
				return ""
			}
			return nation.Country[0].CountryId
		}(),
	}, nil
}

func CreateAgifiedPeople(db *db.DbPeopleRepo, name kafka.PeopleName, info AgifyInfo) error {
	return db.Create(dto.CreatePeople{
		FirstName:  *name.FirstName,
		LastName:   *name.LastName,
		Patronymic: name.Patronymic,
		Age:        info.Age,
		Sex:        info.Sex,
		Nation:     info.Nation,
	})
}
