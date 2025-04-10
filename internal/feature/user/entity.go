package user

import "github.com/google/uuid"

type DBUser struct {
	ID      uuid.UUID `db:"id"`
	Email   string    `db:"email"`
	Version uint64    `db:"version"`
	Role    string    `db:"role"`
}
