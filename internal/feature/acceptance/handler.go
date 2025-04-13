package acceptance

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/Sanchir01/golang-avito/pkg/lib/api"
	sl "github.com/Sanchir01/golang-avito/pkg/lib/log"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

//go:generate go run github.com/vektra/mockery/v2@v2.52.2  --name=AccptanceHandlerInterface --output ./mocks
type AccptanceHandlerInterface interface {
	CreateAcceptanceService(ctx context.Context, pvzId uuid.UUID) (*DBAcceptance, error)
	CloseLastAcceptance(ctx context.Context, pvzID uuid.UUID) (*DBAcceptance, error)
}
type Handler struct {
	Service AccptanceHandlerInterface
	Log     *slog.Logger
}

func NewHandler(s AccptanceHandlerInterface, lg *slog.Logger) *Handler {
	return &Handler{
		Service: s,
		Log:     lg,
	}
}

func (h *Handler) CreateAcceptanceHandler(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.create.acceptance"
	log := h.Log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)
	var req RequestCreateAcceptance

	if err := render.DecodeJSON(r.Body, &req); err != nil {
		log.Error("failed to decode request body", slog.Any("err", err))
		render.JSON(w, r, api.Error("Ошибка при валидации данных pvz"))
		return
	}
	log.Info("request body decoded", slog.Any("request", req))
	if err := validator.New().Struct(req); err != nil {
		log.Error("invalid request", sl.Err(err))
		render.JSON(w, r, api.Error("invalid request"))
		return
	}

	accep, err := h.Service.CreateAcceptanceService(r.Context(), req.PvzId)
	if err != nil {
		log.Error("invalid request", sl.Err(err))
		render.JSON(w, r, api.Error("invalid request"))
		return
	}
	render.JSON(w, r, ResponseCreateAcceptace{
		Response: api.OK(),
		Datetime: accep.CreatedAt,
		ID:       accep.ID,
		PvzId:    accep.PvzId,
		Status:   accep.Status,
	})
}

func (h *Handler) CloseLastAcceptanceHandler(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.close.acceptance"
	log := h.Log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)
	id := chi.URLParam(r, "pvzId")

	parsedID, err := uuid.Parse(id)
	if err != nil {
		log.Error("failed delete product", sl.Err(err))
		render.JSON(w, r, api.Error("failed delete product"))
		return
	}

	accep, err := h.Service.CloseLastAcceptance(r.Context(), parsedID)
	if err != nil {
		log.Error("invalid request", sl.Err(err))
		render.JSON(w, r, api.Error("invalid request"))
		return
	}
	render.JSON(w, r, ResponseCloseAcceptance{
		Response: api.OK(),
		Datetime: accep.CreatedAt,
		ID:       accep.ID,
		PvzId:    accep.PvzId,
		Status:   accep.Status,
	})
}
