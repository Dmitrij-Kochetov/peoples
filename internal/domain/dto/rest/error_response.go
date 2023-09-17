package rest

import (
	"github.com/go-chi/render"
	"net/http"
)

type ErrResponse struct {
	Err            error `json:"-"`
	HTTPStatusCode int   `json:"-"`

	StatusText string `json:"status"`
	AppCode    int64  `json:"code,omitempty"`
	ErrorText  string `json:"error,omitempty"`
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

var (
	ErrNotFound = &ErrResponse{
		HTTPStatusCode: http.StatusNotFound,
		StatusText:     "Resource not found",
	}
	ErrInternalServerError = &ErrResponse{
		HTTPStatusCode: http.StatusInternalServerError,
		StatusText:     "Internal server error",
	}
	ErrBadRequest = &ErrResponse{
		HTTPStatusCode: http.StatusBadRequest,
		StatusText:     "Bad request",
	}
)

func ErrUnprocessableEntity(err error) *ErrResponse {
	return &ErrResponse{
		HTTPStatusCode: http.StatusUnprocessableEntity,
		StatusText:     "Unprocessable entity",
		Err:            err,
		ErrorText:      err.Error(),
	}
}
