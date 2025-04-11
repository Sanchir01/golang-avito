package pvz

import (
	"time"

	"github.com/Sanchir01/golang-avito/pkg/lib/api"
	"github.com/google/uuid"
)

type DBPVZ struct {
	ID               uuid.UUID `db:"id"`
	RegistrationDate time.Time `db:"registration_date"`
	City             string    `db:"city"`
}

type RequestCreatePVZ struct {
	RegistrationDate time.Time `json:"registration_date" validate:"required"`
	City             string    `json:"city" validate:"required"`
}

type ResponseCreatePVZ struct {
	api.Response
	PVZ *DBPVZ
}
