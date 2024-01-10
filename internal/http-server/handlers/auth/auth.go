package auth

import (
	"errors"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"strconv"
	"todo-list/internal/database/postgres"
	"todo-list/internal/lib/api/response"
)

type Request struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type Response struct {
	response.Response
}

type CheckUser interface {
	CheckUser(login, password string) (int64, error)
}

func New(log *slog.Logger, checkUser CheckUser) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server/handlers/auth/auth/New"

		log = slog.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("failed to decode request body", err)

			render.JSON(w, r, response.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			log.Error("invalid request", err)

			render.JSON(w, r, response.Error("invalid request"))

			return
		}

		id, err := checkUser.CheckUser(req.Login, req.Password)
		if errors.Is(err, postgres.ErrUserNotFound) {
			log.Info("user not found", "data", req.Login, req.Password)

			render.JSON(w, r, response.Error("not found"))

			return
		}
		if err != nil {
			log.Error("failed to check user existence")

			render.JSON(w, r, response.Error("internal error"))

			return
		}
		log.Info("check user", slog.String("id", strconv.FormatInt(id, 10)))
		responseOK(w, r)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, Response{
		Response: response.OK(),
	})
}
