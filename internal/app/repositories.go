package app

import "github.com/Sanchir01/golang-avito/internal/feature/user"

type Repositories struct {
	UserRepository *user.Repository
}

func NewRepositories(database *Database) *Repositories {
	return &Repositories{
		UserRepository: user.NewRepository(database.PrimaryDB),
	}
}
