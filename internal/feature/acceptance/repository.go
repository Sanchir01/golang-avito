package acceptance

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/Sanchir01/golang-avito/pkg/lib/api"
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

func (r *Repository) CloseLastAcceptanceStatus(ctx context.Context, pvzID uuid.UUID, tx pgx.Tx) (*DBAcceptance, error) {
	query := `
  WITH last_acceptance AS (
    SELECT id FROM acceptance
    WHERE pvz_id = $1
    ORDER BY created_at DESC
    LIMIT 1
  )
  UPDATE acceptance
  SET status = 'close'
  WHERE id IN (SELECT id FROM last_acceptance)
  RETURNING id, status, created_at, pvz_id
  `

	var updatedAcceptance DBAcceptance

	err := tx.QueryRow(ctx, query, pvzID).Scan(&updatedAcceptance.ID, &updatedAcceptance.Status, &updatedAcceptance.CreatedAt, &updatedAcceptance.PvzId)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, api.ErrNotFoundAcceptance
		}
		return nil, err
	}

	return &updatedAcceptance, nil
}
