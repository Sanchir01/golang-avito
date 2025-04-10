package user

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	primaryDB *pgxpool.Pool
}

func NewRepository(primaryDB *pgxpool.Pool) *Repository {
	return &Repository{primaryDB: primaryDB}
}

func (repo *Repository) Register(ctx context.Context, email, role string, password []byte, tx pgx.Tx) (*DBUser, error) {
	query, arg, err := sq.
		Insert("users").
		Columns("email", "role", "password").
		Values(email, role, password).
		Suffix("RETURNING id, role, email").
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, err
	}
	var users DBUser

	if err := tx.QueryRow(ctx, query, arg...).Scan(&users.ID, &users.Role, &users.Email); err != nil {
		if err == pgx.ErrTxCommitRollback {
			return nil, fmt.Errorf("ошибка при создании пользователя")
		}
		return nil, err
	}

	return &DBUser{
		ID:      users.ID,
		Role:    users.Role,
		Email:   users.Email,
		Version: users.Version,
	}, nil
}
