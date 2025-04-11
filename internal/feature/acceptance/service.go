package acceptance

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ServiceAcceptance interface {
	CreateAcceptance(ctx context.Context, pvzId uuid.UUID, tx pgx.Tx) (*DBAcceptance, error)
}

type Service struct {
	repository ServiceAcceptance
	primaryDB  *pgxpool.Pool
}

func NewService(repository ServiceAcceptance, db *pgxpool.Pool) *Service {
	return &Service{
		repository: repository,
		primaryDB:  db,
	}
}

func (s *Service) CreateAcceptanceService(ctx context.Context, pvzId uuid.UUID) (*DBAcceptance, error) {
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

	accep, err := s.repository.CreateAcceptance(ctx, pvzId, tx)
	if err != nil {
		return nil, err
	}
	return accep, nil
}
