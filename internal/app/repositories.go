package app

import (
	"github.com/Sanchir01/golang-avito/internal/feature/pvz"
	"github.com/Sanchir01/golang-avito/internal/feature/user"
)

type Repositories struct {
	UserRepository *user.Repository
	PVZRepository  *pvz.Repository
}

func NewRepositories(database *Database) *Repositories {
	return &Repositories{
		UserRepository: user.NewRepository(database.PrimaryDB),
		PVZRepository:  pvz.NewRepository(database.PrimaryDB),
	}
}
