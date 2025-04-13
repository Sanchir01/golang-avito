package product

import (
	"time"

	"github.com/Sanchir01/golang-avito/pkg/lib/api"
	"github.com/google/uuid"
)

type DBProduct struct {
	ID          uuid.UUID `db:"id"`
	CreatedAt   time.Time `db:"created_at"`
	Type        string    `db:"type"`
	ReceptionID uuid.UUID `db:"reception_id"`
}

type RequestCreateProduct struct {
	Type        string    `json:"type"`
	AcceptionID uuid.UUID `json:"acception_id"`
}

type ResponseCreateProduct struct {
	api.Response
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	Type        string    `json:"type"`
	AcceptionID uuid.UUID `json:"acception_id"`
}

type RequestDeleteLastProduct struct {
	ReceptionID uuid.UUID `json:"reception_id"`
}

type ResponseDeleteLastProduct struct {
	api.Response
}
