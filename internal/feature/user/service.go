package user

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ServiceUser interface {
	Register(ctx context.Context, email, role string, password []byte, tx pgx.Tx) (*DBUser, error)
}

type Service struct {
	repository ServiceUser
	primaryDB  *pgxpool.Pool
}

func NewService(repository ServiceUser, db *pgxpool.Pool) *Service {
	return &Service{
		repository: repository,
		primaryDB:  db,
	}
}

func (s *Service) RegistrationsService(ctx context.Context, email, role string, password string) (*DBUser, error) {
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
	passwordHash, err := GeneratePasswordHash(password)
	if err != nil {
		return nil, err
	}

	newUser, err := s.repository.Register(ctx, email, role, passwordHash, tx)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}
