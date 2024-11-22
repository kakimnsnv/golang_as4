package auth_repo

import (
	"as4/internal/entity"
	auth_interface "as4/internal/usecase/interface"
	"as4/pkg/postgres"
	"context"

	"go.uber.org/zap"
)

type AuthRepoPostgresImpl struct {
	*postgres.Postgres
	l *zap.Logger
}

var _ auth_interface.AuthRepo = (*AuthRepoPostgresImpl)(nil)

func NewAuthRepoPostgresImpl(p *postgres.Postgres, l *zap.Logger) *AuthRepoPostgresImpl {
	return &AuthRepoPostgresImpl{p, l}
}

func (r *AuthRepoPostgresImpl) CreateUser(ctx context.Context, user entity.User) (entity.User, error) {
	sql, args, err := r.Builder.
		Insert("users").Columns("username", "email", "password", "is_admin").
		Values(user.Username, user.Email, user.Password, user.IsAdmin).
		Suffix("RETURNING id, username, email, is_admin").
		ToSql()
	if err != nil {
		r.l.Error("Failed to build sql", zap.Error(err))
		return entity.User{}, err
	}

	dbUser := entity.User{}

	row := r.Pool.QueryRow(ctx, sql, args...)
	if err = row.Scan(&dbUser.ID, &dbUser.Username, &dbUser.Email, &dbUser.IsAdmin); err != nil {
		r.l.Error("Failed to scan row", zap.Error(err))
		return entity.User{}, err
	}

	return dbUser, nil
}

func (r *AuthRepoPostgresImpl) GetUserByEmail(ctx context.Context, email string) (entity.User, error) {
	sql, args, err := r.Builder.
		Select("id", "username", "email", "password", "is_admin").
		From("users").
		Where("email = ?", email).
		ToSql()
	if err != nil {
		r.l.Error("Failed to build sql", zap.Error(err))
		return entity.User{}, err
	}

	dbUser := entity.User{}

	row := r.Pool.QueryRow(ctx, sql, args...)
	if err = row.Scan(&dbUser.ID, &dbUser.Username, &dbUser.Email, &dbUser.Password, &dbUser.IsAdmin); err != nil {
		r.l.Error("Failed to scan row", zap.Error(err))
		return entity.User{}, err
	}
	return dbUser, err
}
