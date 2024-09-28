package database

import (
	"context"
	"io"

	"github.com/h3th-IV/climateer/pkg/model"
)

type Database interface {
	io.Closer
	CreateUser(ctx context.Context, first_name string, last_name string, email string, password string, phone string, eduInstitute string, sessionkey string) (bool, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	CheckUser(ctx context.Context, email string, password string) (*model.User, error)
	GetBySessionKey(ctx context.Context, sessionkey string) (*model.User, error)
}
