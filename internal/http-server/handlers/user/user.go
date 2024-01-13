package user

import (
	"errors"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"todo-list/internal/database/postgres"
	"todo-list/internal/lib/api/response"
	"todo-list/internal/lib/slogErr"
)

type Request struct {
	Name     string `json:"name" validate:"required"`
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type Response struct {
	response.Response
	Name string `json:"name,omitempty"`
}

type CreateUser interface {
	CreateUser(name, login, password string) (int64, error)
}

func New(log *slog.Logger, user CreateUser) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers/user/user/New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("failed to decode request body", err)
			render.JSON(w, r, response.Error("failed to decode request body"))
			return
		}

		log.Info("request body decoded")

		if err := validator.New().Struct(req); err != nil {
			log.Error("invalid request", err)
			render.JSON(w, r, response.Error("invalid request"))
			return
		}

		id, err := user.CreateUser(req.Name, req.Login, req.Password)
		if errors.Is(err, postgres.ErrUserNameExists) {
			log.Info("username already exists", slog.String("name", req.Name))

			render.JSON(w, r, response.Error("username already exists"))

			return
		} else if errors.Is(err, postgres.ErrUserLoginExists) {
			log.Info("user with this login already exists", slog.String("login", req.Login))

			render.JSON(w, r, response.Error("user with this login already exists"))

			return
		}

		if err != nil {
			log.Error("failed to add user", slogErr.Err(err))

			render.JSON(w, r, response.Error("failed to add user"))

			return
		}

		log.Info("user added", slog.Int64("id", id))
		render.JSON(w, r, Response{Response: response.OK()})
	}
}
