package model

import (
	"context"
	"time"
)

type User struct {
	ID           int       `json:"id"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Email        string    `json:"email"`
	Password     string    `json:"password"`
	Phone        string    `json:"phone"`
	EduInstitute string    `json:"edu_institute"`
	SessionKey   string    `json:"session_key"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Registration struct {
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	Phone        string `json:"phone"`
	EduInstitute string `json:"edu_institute"`
}

type LoginCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// key is an unexported type for keys defined in this package.
// This prevents collisions with keys defined in other packages.
type key struct{}

// playerInfo is the key for model.PlayerInfo values in Contexts. It is
// unexported; clients use model.NewContext and model.FromContext
// instead of using this key directly.
var user key

// NewContext returns a new Context that carries value playerInfo.
func NewContext(ctx context.Context, pi *User) context.Context {
	return context.WithValue(ctx, user, pi)
}

// FromContext returns the User value stored in ctx, if any.
func FromContext(ctx context.Context) (*User, bool) {
	pi, ok := ctx.Value(user).(*User)
	return pi, ok
}
