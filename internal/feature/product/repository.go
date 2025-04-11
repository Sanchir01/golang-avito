package product

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	primaryDB *pgxpool.Pool
}

func NewRepository(primaryDB *pgxpool.Pool) *Repository {
	return &Repository{primaryDB: primaryDB}
}

func (r *Repository) CreateProducts(ctx context.Context, pvzId uuid.UUID, tx pgx.Tx) (*DBProduct, error) {
	query, arg, err := sq.Insert("product").
		Columns("pvz_id").ToSql()
	return nil, nil
}
