package postgres

import (
	"as4/config"
	"context"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Postgres struct {
	maxPoolSize  int
	connAttempts int
	connTimeout  time.Duration

	Builder squirrel.StatementBuilderType
	Pool    *pgxpool.Pool
}

func NewPostgres(lc fx.Lifecycle, config *config.Config, l *zap.Logger) *Postgres {
	pg := &Postgres{
		maxPoolSize:  config.PoolMax,
		connAttempts: 10,
		connTimeout:  time.Second,
	}
	pg.Builder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	poolConfig, err := pgxpool.ParseConfig(config.DSN)
	if err != nil {
		l.Fatal("failed to parse config", zap.Error(err))
	}

	poolConfig.MaxConns = int32(pg.maxPoolSize)

	for pg.connAttempts > 0 {
		pg.Pool, err = pgxpool.ConnectConfig(context.Background(), poolConfig)
		if err == nil {
			l.Info("connected to postgres")
			break
		}

		l.Info("Postgres is trying to connect", zap.Int("attempts", pg.connAttempts))

		time.Sleep(pg.connTimeout)

		pg.connAttempts--
	}

	if err != nil {
		l.Fatal("failed to connect to postgres", zap.Error(err))
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			l.Info("closing postgres connection")
			pg.Close()
			return nil
		},
	})

	return pg
}

func (p *Postgres) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}
