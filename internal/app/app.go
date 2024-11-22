package app

import (
	"as4/config"
	v1 "as4/internal/controller/http/v1"
	auth_usecase "as4/internal/usecase"
	auth_interface "as4/internal/usecase/interface"
	auth_repo "as4/internal/usecase/repo"
	"as4/pkg/httpserver"
	"as4/pkg/postgres"
	"as4/pkg/validator"
	"net/http"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

func Run() {
	fx.New(
		fx.Provide(
			zap.NewDevelopment,
			config.NewConfig,
			postgres.NewPostgres,
			httpserver.NewHTTPSServer,
			fx.Annotate(
				v1.NewRouter,
				fx.As(new(http.Handler)),
			),
			fx.Annotate(
				auth_usecase.NewAuthUseCase,
				fx.As(new(auth_interface.Auth)),
			),
			fx.Annotate(
				auth_repo.NewAuthRepoPostgresImpl,
				fx.As(new(auth_interface.AuthRepo)),
			),
		),
		fx.Invoke(func(*http.Server) { validator.NewValidator() }),
	).Run()
}
