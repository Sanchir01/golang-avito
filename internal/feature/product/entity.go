package product

import (
	"time"

	"github.com/google/uuid"
)

type DBProduct struct {
	ID          uuid.UUID `db:"id"`
	CreatedAt   time.Time `db:"created_at"`
	Type        string    `db:"type"`
	ReceptionID uuid.UUID `db:"reception_id"`
}
