package tests

import (
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/Sanchir01/golang-avito/internal/feature/pvz"
	"github.com/Sanchir01/golang-avito/internal/feature/user"
	"github.com/gavv/httpexpect/v2"
)

func TestIntegrationFlow(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   "localhost:8080",
	}

	e := httpexpect.Default(t, u.String())
	token := e.POST("/api/auth/register").WithJSON(user.RequestRegister{
		Email:    "test13453@test.com",
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

	e.POST("/api/pvz").
		WithHeader("Authorization", "Bearer "+token).
		WithJSON(pvz.RequestCreatePVZ{
			RegistrationDate: time.Now(),
			City:             "Москва",
		}).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object().
		ContainsKey("PVZ").
		ContainsKey("status").
		ValueEqual("status", "OK")
}
