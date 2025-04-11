package user

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/Sanchir01/golang-avito/pkg/lib/api"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ServiceUser interface {
	Register(ctx context.Context, email, role string, password []byte, tx pgx.Tx) (*DBUser, error)
	GetUserByEmail(ctx context.Context, email string) (*DBUser, error)
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

func (s *Service) RegistrationsService(ctx context.Context, email, role, password string) (*DBUser, string, error) {
	conn, err := s.primaryDB.Acquire(ctx)
	if err != nil {
		return nil, "", err
	}
	defer conn.Release()
	tx, err := conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, "", err
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
		return nil, "", err
	}

	newUser, err := s.repository.Register(ctx, email, role, passwordHash, tx)
	if err != nil {
		return nil, "", err
	}

	token, err := GenerateJwtToken(newUser.ID, newUser.Role, newUser.Email, time.Now().Add(14*24*time.Hour))
	if err != nil {
		return nil, "", err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, "", err
	}
	return newUser, token, nil
}

func (s *Service) LoginService(ctx context.Context, email string, password string) (string, error) {
	oldUser, err := s.repository.GetUserByEmail(ctx, email)
	if err != nil {
		return "", err
	}
	verifypass := VerifyPassword(
		oldUser.Password,
		password,
	)
	slog.Error("old user login", oldUser.Role)
	if verifypass == false {
		return "", api.InvalidPassword
	}
	token, err := GenerateJwtToken(oldUser.ID, oldUser.Role, oldUser.Email, time.Now().Add(14*24*time.Hour))
	if err != nil {
		return "", err
	}
	return token, nil
}
