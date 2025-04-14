package product

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

//go:generate go run github.com/vektra/mockery/v2@v2.52.2  --name=HandlerProducts --output ./mocks
type HandlerProducts interface {
	CreateProduct(ctx context.Context, acceptionID uuid.UUID, typeProduct string) (*DBProduct, error)
	DeleteLastProductService(ctx context.Context, AcceptionID uuid.UUID) error
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

// @Summary Создания одно продукта
// @Security ApiKeyAuth
// @Tags product
// @Description Создания одно продукта только для сотрудников
// @Accept json
// @Produce json
// @Param request body RequestCreateProduct true "Данные продукта"
// @Success 200 {object} ResponseCreateProduct
// @Failure 400 {object} api.Response
// @Failure 409 {object} api.Response
// @Failure 500 {object} api.Response
// @Router /api/products [post]
func (h *Handler) CreateProductHandler(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.create.product"
	log := h.Log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)
	var req RequestCreateProduct

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

	product, err := h.Service.CreateProduct(r.Context(), req.AcceptionID, req.Type)
	if err != nil {
		log.Error("invalid request", sl.Err(err))
		render.JSON(w, r, api.Error("invalid request"))
		return
	}

	render.JSON(
		w, r, ResponseCreateProduct{
			Response:    api.OK(),
			ID:          product.ID,
			Type:        product.Type,
			CreatedAt:   product.CreatedAt,
			AcceptionID: product.ReceptionID,
		},
	)
}

// @Summary Удаление последнего добавленного продукта
// @Security ApiKeyAuth
// @Tags product
// @Description Удаление последнего добавленного продукта только для сотрудников
// @Param acceptanceID  path string true "acceptance id"
// @Accept json
// @Produce json
// @Success 200 {object} ResponseCreateProduct
// @Failure 400 {object} api.Response
// @Failure 409 {object} api.Response
// @Failure 500 {object} api.Response
// @Router /api/products/{acceptanceID}/delete_last_product [post]
func (h *Handler) DeleteProductHandler(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.delete.product"
	log := h.Log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)
	id := chi.URLParam(r, "acceptanceID")

	parsedID, err := uuid.Parse(id)
	if err != nil {
		log.Error("failed delete product", sl.Err(err))
		render.JSON(w, r, api.Error("failed delete product"))
		return
	}
	if err := h.Service.DeleteLastProductService(r.Context(), parsedID); err != nil {
		log.Error("failed delete product", sl.Err(err))
		render.JSON(w, r, api.Error("failed delete product"))
		return
	}
	log.Info("success delete product")

	render.JSON(w, r, ResponseDeleteLastProduct{
		Response: api.OK(),
	})
}
