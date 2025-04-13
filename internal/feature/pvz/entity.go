package pvz

import (
	"time"

	pvz "github.com/Sanchir01/golang-avito-proto/pkg/gen/golang/pvz"
	"github.com/Sanchir01/golang-avito/internal/feature/acceptance"
	"github.com/Sanchir01/golang-avito/internal/feature/product"
	"github.com/Sanchir01/golang-avito/pkg/lib/api"
	"github.com/google/uuid"
)

type DBPVZ struct {
	ID               uuid.UUID `db:"id"`
	RegistrationDate time.Time `db:"registration_date"`
	City             string    `db:"city"`
}

type DBReceptionWithProducts struct {
	Reception acceptance.DBAcceptance `json:"reception"`
	Products  []product.DBProduct     `json:"products"`
}

type DBPVZWithReceptions struct {
	PVZ        DBPVZ                     `json:"pvz"`
	Receptions []DBReceptionWithProducts `json:"receptions"`
}

type OnePVZ struct {
	ID               uuid.UUID `db:"id"`
	RegistrationDate time.Time `db:"registration_date"`
	City             string    `db:"city"`
}
type RequestCreatePVZ struct {
	RegistrationDate time.Time `json:"registration_date" validate:"required"`
	City             string    `json:"city" validate:"required"`
}
type ResponseGetAllPVZ struct {
	api.Response
	PVZ []*DBPVZWithReceptions
}
type ResponseGetAllPVZGRPC struct {
	api.Response
	PVZ *pvz.GetPVZListResponse
}

type ResponseCreatePVZ struct {
	api.Response
	PVZ *DBPVZ
}
