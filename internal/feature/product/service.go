package product

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ServiceProduct interface {
	CreateProducts(ctx context.Context, pvzId uuid.UUID, typeProduct string, tx pgx.Tx) (*DBProduct, error)
	DeleteLastProduct(ctx context.Context, AcceptionID uuid.UUID, tx pgx.Tx) error
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

func (s *Service) CreateProduct(ctx context.Context, acceptionID uuid.UUID, typeProduct string) (*DBProduct, error) {
	conn, err := s.primaryDB.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	tx, err := conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback(ctx)
			if rollbackErr != nil {
				err = errors.Join(err, rollbackErr)
				return
			}
		}
	}()

	product, err := s.repository.CreateProducts(ctx, acceptionID, typeProduct, tx)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return product, nil
}

func (s *Service) DeleteLastProductService(ctx context.Context, AcceptionID uuid.UUID) error {
	conn, err := s.primaryDB.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()
	tx, err := conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback(ctx)
			if rollbackErr != nil {
				err = errors.Join(err, rollbackErr)
				return
			}
		}
	}()

	err = s.repository.DeleteLastProduct(ctx, AcceptionID, tx)
	if err != nil {
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}
