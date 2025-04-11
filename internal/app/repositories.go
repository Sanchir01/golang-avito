package app

import (
	"github.com/Sanchir01/golang-avito/internal/feature/acceptance"
	"github.com/Sanchir01/golang-avito/internal/feature/product"
	"github.com/Sanchir01/golang-avito/internal/feature/pvz"
	"github.com/Sanchir01/golang-avito/internal/feature/user"
)

type Repositories struct {
	UserRepository       *user.Repository
	PVZRepository        *pvz.Repository
	AcceptanceRepository *acceptance.Repository
	ProductRepository    *product.Repository
}

func NewRepositories(database *Database) *Repositories {
	return &Repositories{
		UserRepository:       user.NewRepository(database.PrimaryDB),
		PVZRepository:        pvz.NewRepository(database.PrimaryDB),
		AcceptanceRepository: acceptance.NewRepository(database.PrimaryDB),
		ProductRepository:    product.NewRepository(database.PrimaryDB),
	}
}
