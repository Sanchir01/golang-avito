package pvz

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ServicePVZ interface {
	CreatePVZ(ctx context.Context, registerDate time.Time, city string, tx pgx.Tx) (*DBPVZ, error)
	GetAllPVZ(ctx context.Context, startDate, endDate time.Time, page, limit uint64) ([]*DBPVZWithReceptions, error)
}

type Service struct {
	repository ServicePVZ
	primaryDB  *pgxpool.Pool
}

func NewService(repository ServicePVZ, db *pgxpool.Pool) *Service {
	return &Service{
		repository: repository,
		primaryDB:  db,
	}
}

func (s *Service) Create(ctx context.Context, createdDate time.Time, city string) (*DBPVZ, error) {
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

	pvz, err := s.repository.CreatePVZ(ctx, createdDate, city, tx)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return pvz, nil
}

func (s *Service) GetAllPVZService(ctx context.Context, startDate, endDate time.Time, page, limit uint64) ([]*DBPVZWithReceptions, error) {

	pvzs, err := s.repository.GetAllPVZ(ctx, startDate, endDate, page, limit)
	if err != nil {
		return nil, err
	}
	fmt.Println("pvzs", pvzs)
	return pvzs, nil
}
