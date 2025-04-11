package acceptance

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

func (r *Repository) CreateAcceptance(ctx context.Context, pvzId uuid.UUID, tx pgx.Tx) (*DBAcceptance, error) {
	query, arg, err := sq.Insert("acceptance").
		Columns("pvz_id").
		Values(pvzId).
		Suffix("RETURNING id, status, created_at, pvz_id").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	var accep DBAcceptance

	if err := tx.QueryRow(ctx, query, arg...).Scan(&accep.ID, &accep.Status, &accep.CreatedAt, &accep.PvzId); err != nil {
		return nil, err
	}
	return &accep, nil
}
