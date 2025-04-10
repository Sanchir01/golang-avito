package user

import (
	"github.com/Sanchir01/golang-avito/pkg/lib/api"
	"github.com/google/uuid"
)

type DBUser struct {
	ID       uuid.UUID `db:"id"`
	Email    string    `db:"email"`
	Version  uint64    `db:"version"`
	Password []byte    `db:"password"`
	Role     string    `db:"role"`
}

type RequestRegister struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Role     string `json:"role" validate:"required"`
}
type ResponseRegister struct {
	api.Response
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
	Role  string    `json:"role" `
	Token string    `json:"token" `
}

type RequestLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type ResponseLogin struct {
	api.Response
	Token string `json:"token"`
}

type RequestDummyLoggin struct {
	Role string `json:"role" validate:"required"`
}
type ResponseDummyLogin struct {
	api.Response
	Token string `json:"token"`
}
