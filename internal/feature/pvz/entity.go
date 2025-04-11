package pvz

import (
	"time"

	"github.com/google/uuid"
)

type DBPVZ struct {
	id               uuid.UUID `db:"id"`
	registrationDate time.Time `db:"registration_date"`
	city             string    `db:"city"`
}
