package tests

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/Sanchir01/golang-avito/internal/feature/acceptance"
	"github.com/Sanchir01/golang-avito/internal/feature/pvz"
	"github.com/Sanchir01/golang-avito/internal/feature/user"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/gavv/httpexpect/v2"
	"github.com/google/uuid"
)

const host = "localhost:8080"

func TestIntegrationFlow(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
	}

	e := httpexpect.Default(t, u.String())

	token := e.POST("/api/auth/register").WithJSON(user.RequestRegister{
		Email:    gofakeit.Email(),
		Password: "test01",
		Role:     "moderator",
	}).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object().
		Value("token").
		String().
		Raw()

	responsepvz := e.POST("/api/pvz").
		WithHeader("Authorization", "Bearer "+token).
		WithJSON(pvz.RequestCreatePVZ{
			RegistrationDate: time.Now(),
			City:             "Москва",
		}).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object()

	id := responsepvz.Value("PVZ").Object().Value("ID").String().Raw()
	pvzuuid, err := uuid.Parse(id)
	if err != nil {
		t.Error(err.Error())
	}
	tokenemployee := e.POST("/api/auth/register").WithJSON(user.RequestRegister{
		Email:    gofakeit.Email(),
		Password: "test01",
		Role:     "employee",
	}).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object().
		Value("token").
		String().
		Raw()

	acceptanceresponse := e.POST("/api/receptions").
		WithHeader("Authorization", "Bearer "+tokenemployee).
		WithJSON(acceptance.RequestCreateAcceptance{
			PvzId: pvzuuid,
		}).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object()

	fmt.Println("receptions", acceptanceresponse)
}
