package list

import (
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
	UserId int64 `json:"userId" validate:"required"`
}

type Response struct {
	Task []postgres.Task `json:"task,omitempty"`
	response.Response
}

type List interface {
	ListOfTask(userId int64) ([]postgres.Task, error)
}

func New(log *slog.Logger, list List) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		const op = "handlers/list/list/New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		if err := render.DecodeJSON(r.Body, &req); err != nil {
			slog.Error("failed to decode request body: ", err)

			render.JSON(w, r, response.Error("failed to decode request body"))

			return
		}

		log.Info("request body decoded", req)

		if err := validator.New().Struct(req); err != nil {
			log.Error("invalid request", err)

			render.JSON(w, r, response.Error("invalid request"))

			return
		}

		results, err := list.ListOfTask(req.UserId)
		if err != nil {
			log.Error("failed to get list of task ", slogErr.Err(err))

			render.JSON(w, r, response.Error("internal error"))

			return
		}

		log.Info("got list of task")

		render.JSON(w, r, results)
	}
}
