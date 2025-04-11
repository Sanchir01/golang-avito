package custommiddleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	contextkey "github.com/Sanchir01/golang-avito/internal/context"
	"github.com/Sanchir01/golang-avito/internal/feature/user"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	requestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "auth",
			Subsystem: "http",
			Name:      "request_total",
			Help:      "Total number of HTTP requests",
		},
		[]string{"path", "method"},
	)

	requestDuration = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: "auth",
			Subsystem: "http",
			Name:      "request_duration_seconds",
			Help:      "Duration of HTTP requests in seconds",
			Objectives: map[float64]float64{
				0.5:  0.05,
				0.9:  0.01,
				0.99: 0.001,
			},
		},
		[]string{"path", "method"},
	)
)

func AuthMiddleware(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "invalid Authorization format", http.StatusUnauthorized)
				return
			}

			users, err := user.ParseToken(parts[1])
			if err != nil {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			fmt.Println("middleware role", users, allowedRoles)
			if !hasRole(users.Role, allowedRoles) {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			ctx := context.WithValue(r.Context(), contextkey.UserIDCtxKey, users)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func PrometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount.WithLabelValues(r.URL.Path, r.Method).Inc()
		next.ServeHTTP(w, r)
	})
}

func hasRole(userRole string, allowedRoles []string) bool {
	for _, role := range allowedRoles {
		if userRole == role {
			return true
		}
	}
	return false
}
