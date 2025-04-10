package app

import "github.com/Sanchir01/golang-avito/internal/feature/user"

type Services struct {
	UserService *user.Service
}

func NewServices(r *Repositories, db *Database) *Services {
	return &Services{
		UserService: user.NewService(r.UserRepository),
	}
}
