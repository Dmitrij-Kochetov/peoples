package kafka

type PeopleName struct {
	FirstName  *string `json:"name"`
	LastName   *string `json:"surname"`
	Patronymic string  `json:"patronymic,omitempty"`
}
