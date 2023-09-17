package rest

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/Dmitrij-Kochetov/peoples/internal/adapter/http-server/middleware/logger"
	"github.com/Dmitrij-Kochetov/peoples/internal/adapter/logging"
	"github.com/Dmitrij-Kochetov/peoples/internal/application/usecases"
	"github.com/Dmitrij-Kochetov/peoples/internal/domain/dto"
	"github.com/Dmitrij-Kochetov/peoples/internal/domain/dto/rest"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

func (s *Server) setupRoutes() {
	s.router.Use(render.SetContentType(render.ContentTypeJSON))

	s.router.Use(middleware.Recoverer)
	s.router.Use(middleware.RequestID)
	s.router.Use(logger.New(s.logger))

	s.router.Route("/api/v1", func(r chi.Router) {
		r.Route("/peoples", func(r chi.Router) {
			r.Get("/", s.getPeoples)
			r.Get("/{id}", s.getPeople)
			r.Post("/", s.createPeople)
			r.Put("/{id}", s.updatePeople)
			r.Delete("/{id}", s.deletePeople)
		})
	})
}

func (s *Server) handleError(w http.ResponseWriter, r *http.Request, e *rest.ErrResponse) {
	err := render.Render(w, r, e)
	if err != nil {
		s.logger.Error("failed to render error", logging.Err(err))
	}
}

func (s *Server) getPeoples(w http.ResponseWriter, r *http.Request) {
	data := &rest.FilterRequest{}
	if err := render.Bind(r, data); err != nil {
		s.logger.Error("bad request", logging.Err(err))
		s.handleError(w, r, rest.ErrBadRequest)
		return
	}

	peoples, err := usecases.GetAllPeopleByFilter(context.Background(), s.repo, dto.Filter(*data))
	if err != nil {
		s.logger.Error("failed to get peoples", logging.Err(err))
		s.handleError(w, r, rest.ErrInternalServerError)
		return
	}

	if err := render.RenderList(w, r, rest.NewListPeopleResponse(peoples)); err != nil {
		s.logger.Error("failed to render", logging.Err(err))
	}
}

func (s *Server) getPeople(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		s.logger.Error("error parsing", logging.Err(err))
		s.handleError(w, r, rest.ErrBadRequest)
		return
	}

	people, err := usecases.GetPeopleByID(context.Background(), s.repo, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.handleError(w, r, rest.ErrNotFound)
		} else {
			s.logger.Info("error getting people by id", logging.Err(err))
			s.handleError(w, r, rest.ErrInternalServerError)
		}
		return
	}

	pr := rest.NewPeopleResponse(*people)
	if err := render.Render(w, r, pr); err != nil {
		s.logger.Error("error rendering", logging.Err(err))
	}
}

func (s *Server) createPeople(w http.ResponseWriter, r *http.Request) {
	data := &rest.CreatePeopleRequest{}
	if err := render.Bind(r, data); err != nil {
		s.logger.Error("bad request", logging.Err(err))
		s.handleError(w, r, rest.ErrBadRequest)
		return
	}

	if data.Sex != "male" && data.Sex != "female" {
		s.logger.Error("bad request", logging.Err(fmt.Errorf("sex must be male | female")))
		s.handleError(w, r, rest.ErrBadRequest)
		return
	}

	err := usecases.CreatePeople(context.Background(), s.repo, dto.CreatePeople(*data))
	if err != nil {
		s.logger.Error("internal server error", logging.Err(err))
		s.handleError(w, r, rest.ErrInternalServerError)
		return
	}

	w.WriteHeader(200)
	if _, err = w.Write(nil); err != nil {
		s.logger.Error("failed to write response", logging.Err(err))
	}
}

func (s *Server) updatePeople(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		s.logger.Error("error parsing", logging.Err(err))
		s.handleError(w, r, rest.ErrBadRequest)
		return
	}

	data := &rest.CreatePeopleRequest{}
	if err := render.Bind(r, data); err != nil {
		s.logger.Error("bad request", logging.Err(err))
		s.handleError(w, r, rest.ErrBadRequest)
		return
	}

	if data.Sex != "male" && data.Sex != "female" {
		s.logger.Error("bad request", logging.Err(fmt.Errorf("sex must be male | female")))
		s.handleError(w, r, rest.ErrBadRequest)
		return
	}

	if err := usecases.UpdatePeopleByID(context.Background(), s.repo, dto.People{
		ID:         id,
		FirstName:  data.FirstName,
		LastName:   data.LastName,
		Patronymic: data.Patronymic,
		Age:        data.Age,
		Sex:        data.Sex,
		Nation:     data.Nation,
	}); err != nil {
		s.logger.Error("internal serever error", logging.Err(err))
		s.handleError(w, r, rest.ErrInternalServerError)
		return
	}

	if _, err = w.Write(nil); err != nil {
		s.logger.Error("failed to write response", logging.Err(err))
	}
}

func (s *Server) deletePeople(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		s.logger.Error("error parsing", logging.Err(err))
		s.handleError(w, r, rest.ErrBadRequest)
		return
	}

	err = usecases.DeletePeopleByID(context.Background(), s.repo, id)
	if err != nil {
		s.logger.Error("internal serever error", logging.Err(err))
		s.handleError(w, r, rest.ErrInternalServerError)
		return
	}

	if _, err = w.Write(nil); err != nil {
		s.logger.Error("failed to write response", logging.Err(err))
	}
}
