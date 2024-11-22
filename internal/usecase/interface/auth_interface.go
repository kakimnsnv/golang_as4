package auth_interface

import (
	"as4/config"
	"as4/internal/entity"
	"context"
)

type (
	Auth interface {
		Login(ctx context.Context, email, password string, config *config.Config) (entity.AuthResponse, error)
		Register(ctx context.Context, email, password, username string, config *config.Config) (entity.AuthResponse, error)
	}

	AuthRepo interface {
		CreateUser(ctx context.Context, user entity.User) (entity.User, error)
		GetUserByEmail(ctx context.Context, email string) (entity.User, error)
	}
)
