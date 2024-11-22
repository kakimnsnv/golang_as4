package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"go.uber.org/zap"
)

type (
	Config struct {
		HTTP `yaml:"http"`
		App  `yaml:"app"`
		PG   `yaml:"postgres"`
		Log  `yaml:"log"`
		Auth `yaml:"auth"`
	}

	App struct {
		Name    string `env-required:"true" yaml:"name" env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	HTTP struct {
		Port string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
	}

	PG struct {
		PoolMax int    `env-required:"true" yaml:"pool_max" env:"PG_POOL_MAX"`
		DSN     string `env-required:"true" yaml:"dsn" env:"PG_DSN"`
	}

	Log struct {
		Level string `env-required:"true" yaml:"level" env:"LOG_LEVEL"`
	}

	Auth struct {
		JWTSecret  string `env-required:"true" yaml:"jwt" env:"JWT_SECRET"`
		CSRFSecret string `env-required:"true" yaml:"csrf" env:"CSRF_SECRET"`
	}
)

func NewConfig(l *zap.Logger) *Config {
	cfg := &Config{}

	err := cleanenv.ReadConfig("./config/config.yml", cfg)
	if err != nil {
		l.Fatal("Failed to read config", zap.Error(err))
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		l.Fatal("Failed to read env", zap.Error(err))
	}

	return cfg
}
