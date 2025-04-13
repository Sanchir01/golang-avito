package product

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

func (r *Repository) CreateProducts(ctx context.Context, acceptanceID uuid.UUID, typeProduct string, tx pgx.Tx) (*DBProduct, error) {
	query, arg, err := sq.
		Insert("product").
		Columns("receiving_id", "type").
		Values(acceptanceID, typeProduct).
		Suffix("RETURNING id, type, receiving_id,created_at").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}
	var product DBProduct

	if err := tx.QueryRow(ctx, query, arg...).Scan(&product.ID, &product.Type, &product.ReceptionID, &product.CreatedAt); err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *Repository) DeleteLastProduct(ctx context.Context, AcceptionID uuid.UUID, tx pgx.Tx) error {
	query := `
        WITH last_product AS (
            SELECT id FROM product
            WHERE receiving_id = $1
            ORDER BY created_at DESC
            LIMIT 1
        )
        DELETE FROM product
        WHERE id IN (SELECT id FROM last_product)
    `

	_, err := tx.Exec(ctx, query, AcceptionID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return api.ErrNotFoundProduct
		}
		return err
	}
	return nil
}
