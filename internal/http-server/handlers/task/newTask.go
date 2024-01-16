package task

import (
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"time"
	"todo-list/internal/lib/api/response"
	"todo-list/internal/lib/slogErr"
)

type Request struct {
	Name        string    `json:"name" validate:"required"`
	Description string    `json:"description,omitempty"`
	Deadline    time.Time `json:"deadline" validate:"required"`
	UserID      int64     `json:"userID" validate:"required"`
}

type Response struct {
	response.Response
	Name string `json:"name,omitempty"`
}

type CreateTask interface {
	CreateTask(userId int64, name, description string, date time.Time) (int64, error)
}

func New(log *slog.Logger, createTask CreateTask) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers/task/newTask/New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("failed to decode request body", err)

			render.JSON(w, r, response.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded: ", req)

		if err := validator.New().Struct(req); err != nil {
			log.Error("invalid request", err)

			render.JSON(w, r, "invalid request")

			return
		}

		id, err := createTask.CreateTask(req.UserID, req.Name, req.Description, req.Deadline)
		if err != nil {
			log.Error("failed to add task", slogErr.Err(err))

			render.JSON(w, r, response.Error("failed to add task"))

			return
		}

		log.Info("created task: ", id)

		render.JSON(w, r, Response{Response: response.OK(), Name: req.Name})
	}
}
