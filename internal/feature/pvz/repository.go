package pvz

import (
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	primaryDB *pgxpool.Pool
}

func NewRepository(primaryDB *pgxpool.Pool) *Repository {
	return &Repository{primaryDB: primaryDB}
}

func (h *Handler) GetAllPVZ(startDate, endDate time.Time, page, limit int32)
