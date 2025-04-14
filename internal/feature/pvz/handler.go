package pvz

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	grpcserver "github.com/Sanchir01/golang-avito/internal/server/servers/grpc"
	"github.com/Sanchir01/golang-avito/pkg/lib/api"
	sl "github.com/Sanchir01/golang-avito/pkg/lib/log"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

//go:generate go run github.com/vektra/mockery/v2@v2.52.2 --name=HandlersInterface --output ./mocks
type HandlersInterface interface {
	Create(ctx context.Context, createdDate time.Time, city string) (*DBPVZ, error)
	GetAllPVZService(ctx context.Context, startDate, endDate time.Time, page, limit uint64) ([]*DBPVZWithReceptions, error)
}

type Handler struct {
	Service HandlersInterface
	GRPCPVZ *grpcserver.GRPCClientPVZ
	Log     *slog.Logger
}

func NewHandler(s HandlersInterface, lg *slog.Logger, pvzgrpc *grpcserver.GRPCClientPVZ) *Handler {
	return &Handler{
		Service: s,
		Log:     lg,
		GRPCPVZ: pvzgrpc,
	}
}

// @Summary Создание пункта выдачи
// @Security ApiKeyAuth
// @Tags pvz
// @Description Создание пункта выдачи только для администраторов
// @Accept json
// @Produce json
// @Param request body RequestCreatePVZ true "Данные тестового логина"
// @Success 200 {object} ResponseCreatePVZ
// @Failure 400 {object} api.Response
// @Failure 409 {object} api.Response
// @Failure 500 {object} api.Response
// @Router /api/pvz [post]
func (h *Handler) CreatePVZHandler(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.create.pvz"
	log := h.Log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)
	var req RequestCreatePVZ

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

	pvz, err := h.Service.Create(r.Context(), req.RegistrationDate, req.City)
	if err != nil {
		log.Error("invalid request", sl.Err(err))
		render.JSON(w, r, api.Error("invalid request"))
		return
	}

	render.JSON(
		w, r, ResponseCreatePVZ{
			Response: api.OK(),
			PVZ:      pvz,
		})
}

// @Summary Получение всех пвз со всеми приемками товаров и всеми товарами этих приемок
// @Security ApiKeyAuth
// @Tags pvz
// @Description Получение всех пвз со всеми приемками товаров и всеми товарами этих приемок администраторов и сотрудников, надо включить grpc server
// @Accept json
// @Produce json
// @Success 200 {object} ResponseGetAllPVZ
// @Failure 400 {object} api.Response
// @Failure 409 {object} api.Response
// @Failure 500 {object} api.Response
// @Router /api/pvz [get]
func (h *Handler) GetAllPVZHandler(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.get.all.pvz"
	log := h.Log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	var startDate, endDate time.Time
	var page, limit uint64 = 1, 20

	if startDateStr != "" {
		parsedStartDate, err := time.Parse(time.RFC3339, startDateStr)
		if err != nil {
			render.JSON(w, r, api.Error("invalid request"))
			return
		}
		startDate = parsedStartDate
	} else {
		startDate = time.Now().AddDate(0, 0, -7)
	}

	if endDateStr != "" {
		parsedEndDate, err := time.Parse(time.RFC3339, endDateStr)
		if err != nil {
			render.JSON(w, r, api.Error("invalid request"))
			return
		}
		endDate = parsedEndDate
	}

	if pageStr != "" {
		parsedPage, err := strconv.ParseUint(pageStr, 10, 64)
		if err != nil || parsedPage < 1 {
			render.JSON(w, r, api.Error("invalid request"))
			return
		}
		page = parsedPage
	}

	if limitStr != "" {
		parsedLimit, err := strconv.ParseUint(limitStr, 10, 64)
		if err != nil || parsedLimit < 1 {
			render.JSON(w, r, api.Error("invalid request"))
			return
		}
		limit = parsedLimit
	}

	pvzs, err := h.Service.GetAllPVZService(r.Context(), startDate, endDate, page, limit)
	if err != nil {
		log.Error("failed to get all pvz", sl.Err(err))
		render.JSON(w, r, api.Error("failed to get all pvz"))
		return
	}

	render.JSON(w, r, ResponseGetAllPVZ{
		Response: api.OK(),
		PVZ:      pvzs,
	})
}

// @Summary Получение всех пвз
// @Security ApiKeyAuth
// @Tags pvz
// @Description Получение всех пвз только для администраторов и сотрудников, надо включить grpc server
// @Accept json
// @Produce json
// @Success 200 {object} ResponseGetAllPVZ
// @Failure 400 {object} api.Response
// @Failure 409 {object} api.Response
// @Failure 500 {object} api.Response
// @Router /api/pvz_grpc [get]
func (h *Handler) GetAllGRPCPVZHandler(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.getall.pvz.grpc"
	log := h.Log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	pvzs, err := h.GRPCPVZ.AllPVZHandler(r.Context())
	if err != nil {
		log.Error("failed to get all pvz", sl.Err(err))
		render.JSON(w, r, api.Error("failed to get all pvz"))
		return
	}

	render.JSON(w, r, ResponseGetAllPVZGRPC{
		Response: api.OK(),
		PVZ:      pvzs,
	})
}
