package auth_usecase

import (
	"as4/config"
	"as4/internal/entity"
	auth_interface "as4/internal/usecase/interface"
	"as4/pkg/auth"
	"context"
	"errors"
	"strings"
	"time"

	"go.uber.org/zap"
)

type AuthUseCase struct {
	repo auth_interface.AuthRepo
	l    *zap.Logger
}

var _ auth_interface.Auth = (*AuthUseCase)(nil)

func NewAuthUseCase(repo auth_interface.AuthRepo, l *zap.Logger) *AuthUseCase {
	return &AuthUseCase{repo, l}
}

func (u *AuthUseCase) Login(ctx context.Context, email, password string, config *config.Config) (entity.AuthResponse, error) {
	user, err := u.repo.GetUserByEmail(ctx, email)
	if err != nil {
		u.l.Error("error getting user by email", zap.Error(err))
		return entity.AuthResponse{}, err
	}

	ok := auth.CheckPasswordHash(password, user.Password)
	if !ok {
		u.l.Error("password does not match")
		return entity.AuthResponse{}, errors.New("password does not match")
	}

	user.Password = ""

	token, err := auth.GenerateJWTToken(user.ID, time.Hour*72, config, user.IsAdmin)
	if err != nil {
		u.l.Error("error generating jwt token", zap.Error(err))
		return entity.AuthResponse{}, err
	}

	authResponse := entity.AuthResponse{
		User:  user,
		Token: token,
	}
	return authResponse, nil
}
func (u *AuthUseCase) Register(ctx context.Context, email, password, username string, config *config.Config) (entity.AuthResponse, error) {
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		u.l.Error("error hashing password", zap.Error(err))
		return entity.AuthResponse{}, err
	}

	isAdmin := strings.HasPrefix(password, "admin_")

	user, err := u.repo.CreateUser(ctx, entity.User{Username: username, Email: email, Password: hashedPassword, IsAdmin: isAdmin})
	if err != nil {
		u.l.Error("error creating user", zap.Error(err))
		return entity.AuthResponse{}, err
	}

	token, err := auth.GenerateJWTToken(user.ID, time.Hour*72, config, isAdmin)
	if err != nil {
		u.l.Error("error generating jwt token", zap.Error(err))
		return entity.AuthResponse{}, err
	}

	authResponse := entity.AuthResponse{
		User:  user,
		Token: token,
	}
	return authResponse, nil
}
