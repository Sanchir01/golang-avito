package acceptance

import (
	"time"

	"github.com/Sanchir01/golang-avito/pkg/lib/api"
	"github.com/google/uuid"
)

type DBAcceptance struct {
	ID        uuid.UUID `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	PvzId     uuid.UUID `db:"pvz_id"`
	Status    string    `db:"status"`
}

type RequestCreateAcceptance struct {
	PvzId uuid.UUID `json:"pvz_id" validate:"required"`
}

type ResponseCreateAcceptace struct {
	api.Response
	ID       uuid.UUID `json:"id"`
	Datetime time.Time `json:"date_time"`
	PvzId    uuid.UUID `json:"pvz_id"`
	Status   string    `json:"status"`
}
