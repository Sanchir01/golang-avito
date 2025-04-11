package app

import (
	"github.com/Sanchir01/golang-avito/internal/feature/acceptance"
	"github.com/Sanchir01/golang-avito/internal/feature/pvz"
	"github.com/Sanchir01/golang-avito/internal/feature/user"
)

type Services struct {
	UserService       *user.Service
	PVZService        *pvz.Service
	AcceptanceService *acceptance.Service
}

func NewServices(r *Repositories, db *Database) *Services {
	return &Services{
		UserService:       user.NewService(r.UserRepository, db.PrimaryDB),
		PVZService:        pvz.NewService(r.PVZRepository, db.PrimaryDB),
		AcceptanceService: acceptance.NewService(r.AcceptanceRepository, db.PrimaryDB),
	}
}
