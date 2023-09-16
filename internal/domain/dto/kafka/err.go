package kafka

type Error struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

var (
	ErrRequiredFieldNotExists = Error{
		Message: "validation failed",
		Error:   "required field not exists",
	}
	ErrRequiredFieldIsEmpty = Error{
		Message: "validation failed",
		Error:   "required field is empty",
	}
)
