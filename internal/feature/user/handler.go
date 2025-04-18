package user

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/Sanchir01/golang-avito/pkg/lib/api"
	sl "github.com/Sanchir01/golang-avito/pkg/lib/log"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

//go:generate go run github.com/vektra/mockery/v2@v2.52.2 --name=HandlerUser --output ./mocks
type HandlerUser interface {
	LoginService(ctx context.Context, email string, password string) (string, error)
	RegistrationsService(ctx context.Context, email, role, password string) (*DBUser, string, error)
}
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

// @Summary Регистрация пользователя
// @Tags auth
// @Description Регистрация нового пользователя в системе
// @Accept json
// @Produce json
// @Param request body RequestRegister true "Данные для регистрации"
// @Success 200 {object} ResponseRegister
// @Failure 400 {object} api.Response
// @Failure 409 {object} api.Response
// @Failure 500 {object} api.Response
// @Router /api/auth/register [post]
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
	newuser, token, err := h.Service.RegistrationsService(r.Context(), req.Email, req.Role, req.Password)
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

	render.JSON(w, r, ResponseRegister{
		Response: api.OK(),
		ID:       newuser.ID,
		Role:     newuser.Role,
		Email:    newuser.Email,
		Token:    token,
	})
}

// @Summary Авторизация пользователя
// @Tags auth
// @Description Авторизация пользователя в системе
// @Accept json
// @Produce json
// @Param request body RequestLogin true "Данные логина"
// @Success 200 {object} ResponseRegister
// @Failure 400 {object} api.Response
// @Failure 409 {object} api.Response
// @Failure 500 {object} api.Response
// @Router /api/auth/login [post]
func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.login"
	log := h.Log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)
	var req RequestLogin

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

	token, err := h.Service.LoginService(r.Context(), req.Email, req.Password)
	if errors.Is(err, api.InvalidPassword) {
		log.Error("invalid request", sl.Err(err))
		render.JSON(w, r, api.Error(err.Error()))
		return
	}
	if err != nil {
		log.Error("invalid request", sl.Err(err))
		render.JSON(w, r, api.Error(""))
		return
	}
	log.Info("success login")

	render.JSON(w, r, ResponseLogin{
		Response: api.OK(),
		Token:    token,
	})
}

// @Summary Получение токена
// @Tags auth
// @Description Авторизация пользователя в системе
// @Accept json
// @Produce json
// @Param request body RequestLogin true "Данные тестового логина"
// @Success 200 {object} RequestDummyLoggin
// @Failure 400 {object} api.Response
// @Failure 409 {object} api.Response
// @Failure 500 {object} api.Response
// @Router /api/auth/dummyLogin [post]
func (h *Handler) DummyLoginHandler(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.login"
	log := h.Log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)
	var req RequestDummyLoggin

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
	token, err := GenerateJwtToken(uuid.New(), req.Role, "", time.Now().Add(14*24*time.Hour))
	if err != nil {
		log.Error("invalid request", sl.Err(err))
		render.JSON(w, r, api.Error("invalid request"))
		return
	}

	log.Info("success dummy login")

	render.JSON(w, r, ResponseDummyLogin{
		Response: api.OK(),
		Token:    token,
	})
}
