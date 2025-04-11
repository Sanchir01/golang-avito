package acceptance

import "github.com/jackc/pgx/v5/pgxpool"

type ServiceAcceptance interface {
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
