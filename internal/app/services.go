package app

import ()

type Services struct {
}

func NewServices(r *Repositories, db *Database) *Services {
	return &Services{}
}
