package auth

import (
	"errors"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"log/slog"
	"net/http"
	"time"
	"todo-list/internal/database/postgres"
	"todo-list/internal/lib/api/response"
	"todo-list/internal/lib/slogErr"
)

type Request struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type Response struct {
	Token string `json:"token,omitempty"`
	response.Response
}

type CheckUser interface {
	CheckUser(login, password string) (int64, error)
}

type TokenClaims struct {
	jwt.StandardClaims
	UserId int64 `json:"user_id"`
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
			log.Info("signIn not found", "data", req.Login, req.Password)

			render.JSON(w, r, response.Error("not found"))

			return
		}
		if err != nil {
			log.Error("failed to check signIn existence", slogErr.Err(err))

			render.JSON(w, r, response.Error("internal error"))

			return
		}

		log.Info("check signIn", slog.Int64("id", id))

		token, err := GenerateToken(id)
		if err != nil {
			log.Error("failed to create token", slogErr.Err(err))

			render.JSON(w, r, response.Error("internal error"))

			return
		}

		render.JSON(w, r, Response{Token: token, Response: response.OK()})

		http.Redirect(w, r, "/list", http.StatusSeeOther)
	}
}

const signingKey = "fhdsjalhhgFduighhuGHHUI321hfuDSAHO"

func GenerateToken(userId int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaims{jwt.StandardClaims{
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		IssuedAt:  time.Now().Unix(),
	},
		userId,
	})

	return token.SignedString([]byte(signingKey))
}
