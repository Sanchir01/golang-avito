package user

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/Sanchir01/golang-avito/pkg/lib/api"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	primaryDB *pgxpool.Pool
}

func NewRepository(primaryDB *pgxpool.Pool) *Repository {
	return &Repository{primaryDB: primaryDB}
}

func (repo *Repository) Register(ctx context.Context, email, role string, password []byte, tx pgx.Tx) (*DBUser, error) {
	fmt.Println("protp repository", email, role, string(password))
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
		var pgErr *pgconn.PgError
		if err == pgx.ErrTxCommitRollback {
			return nil, fmt.Errorf("ошибка при создании пользователя")
		}
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, api.ErrUserAlreadyExists
		}
		return nil, err
	}

	return &DBUser{
		ID:    users.ID,
		Role:  users.Role,
		Email: users.Email,
	}, nil
}

func (repo *Repository) GetUserByEmail(ctx context.Context, email string) (*DBUser, error) {
	conn, err := repo.primaryDB.Acquire(ctx)
	if err != nil {
		return nil, err
	}

	defer conn.Release()

	query, arg, err := sq.Select("id,password,role").
		From("users").Where(sq.Eq{"email": email}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}
	var user DBUser
	if err := conn.QueryRow(ctx, query, arg...).Scan(&user.ID, &user.Password, &user.Role); err != nil {
		return nil, err
	}
	return &user, nil
}
