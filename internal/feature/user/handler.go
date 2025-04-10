package user

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/Sanchir01/golang-avito/pkg/lib/api"
	sl "github.com/Sanchir01/golang-avito/pkg/lib/log"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	Service *Service
	Log     *slog.Logger
}

func NewHandler(s *Service, lg *slog.Logger) *Handler {
	return &Handler{
		Service: s,
		Log:     lg,
	}
}
func (h *Handler) RegistrationHandler(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.register"
	log := h.Log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)
	var req RequestRegister

	if err := render.DecodeJSON(r.Body, &req); err != nil {
		log.Error("failed to decode request body", slog.Any("err", err))
		render.JSON(w, r, api.Error("Ошибка при валидации данных"))
		return
	}
	log.Info("request body decoded", slog.Any("request", req))
	if err := validator.New().Struct(req); err != nil {
		log.Error("invalid request", sl.Err(err))
		render.JSON(w, r, api.Error("invalid request"))
		return
	}
	newuser, err := h.Service.RegistrationsService(r.Context(), req.Email, req.Role, req.Password)
	if errors.Is(err, api.InvalidPassword) {
		log.Info("password error", sl.Err(err))
		render.JSON(w, r, api.Error("Введен неправильный пароль"))
		return
	}
	if err != nil {
		log.Error("failed auth user", sl.Err(err))
		render.JSON(w, r, api.Error("failed, auth user"))
		return
	}
	log.Info("success register")
	token, err := GenerateJwtToken(newuser.ID, newuser.Role, newuser.Email, time.Now().Add(14*24*time.Hour))

	render.JSON(w, r, ResponseRegister{
		Response: api.OK(),
		ID:       newuser.ID,
		Role:     newuser.Role,
		Email:    newuser.Email,
		Token:    token,
	})
}
