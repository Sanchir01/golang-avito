package product

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ServiceProduct interface {
	CreateProducts(ctx context.Context, pvzId uuid.UUID, tx pgx.Tx) (*DBProduct, error)
}

type Service struct {
	repository ServiceProduct
	primaryDB  *pgxpool.Pool
}

func NewService(repository ServiceProduct, db *pgxpool.Pool) *Service {
	return &Service{
		repository: repository,
		primaryDB:  db,
	}
}
