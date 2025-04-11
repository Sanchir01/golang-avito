package pvz

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/Sanchir01/golang-avito/pkg/lib/api"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	primaryDB *pgxpool.Pool
}

func NewRepository(primaryDB *pgxpool.Pool) *Repository {
	return &Repository{primaryDB: primaryDB}
}

//func (r *Repository) GetAllPVZ(ctx context.Context, startDate, endDate time.Time, page, limit uint64) ([]*DBPVZ, error) {
//	conn, err := r.primaryDB.Acquire(ctx)
//	if err != nil {
//		return nil, api.ErrorCreateQueryString
//	}
//
//	query, arg, err := sq.Select("id, registration_date, city").
//		From("pvz").
//		Where("").
//		Limit(limit).
//		PlaceholderFormat(sq.Dollar).
//		ToSql()
//	if err != nil {
//		return nil, api.ErrorCreateQueryString
//	}
//
//	return nil, nil
//}

func (r *Repository) CreatePVZ(ctx context.Context, registerDate time.Time, city string, tx pgx.Tx) (*DBPVZ, error) {
	query, arg, err := sq.Insert("pvz").
		Columns("registration_date", "city").
		Values(registerDate, city).
		Suffix("RETURNING id, registration_date, city").
		PlaceholderFormat(sq.Dollar).
		ToSql()

	var pvz DBPVZ
	if err := tx.QueryRow(ctx, query, arg...).Scan(&pvz.ID, &pvz.RegistrationDate, &pvz.City); err != nil {
		if err == pgx.ErrTxCommitRollback {
			return nil, api.ErrCreatePvz
		}
		return nil, err
	}
	if err != nil {
		return nil, api.ErrorCreateQueryString
	}

	return &DBPVZ{
		ID:               pvz.ID,
		RegistrationDate: pvz.RegistrationDate,
		City:             pvz.City,
	}, nil
}
