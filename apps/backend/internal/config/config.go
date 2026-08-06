package config

import (
	"github.com/go-playground/validator/v10"
	_ "github.com/joho/godotenv/autoload"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/v2"
	"github.com/rs/zerolog"
	"os"
	"strings"
)

type Config struct {
	Primary       Primary              `koanf:"primary" validate:"required"`
	Server        ServerConfig         `koanf:"server" validate:"required"`
	Database      DatabaseConfig       `koanf:"database" validate:"required"`
	Auth          AuthConfig           `koanf:"auth" validate:"required"`
	Redis         RedisConfig          `koanf:"redis" validate:"required"`
	Observability *ObservabilityConfig `koanf:"observability"`
	Integration   IntegrationConfig    `koanf:"integration" validate:"required"`
	AWS           AWSConfig            `koanf:"aws"`
	Search        SearchConfig         `koanf:"search"`
}

type Primary struct {
	Env string `koanf:"env" validate:"required"`
}

type ServerConfig struct {
	Port               string   `koanf:"port" validate:"required"`
	ReadTimeout        int      `koanf:"read_timeout" validate:"required"`
	WriteTimeout       int      `koanf:"write_timeout" validate:"required"`
	IdleTimeout        int      `koanf:"idle_timeout" validate:"required"`
	CORSAllowedOrigins []string `koanf:"cors_allowed_origins" validate:"required"`
}

type DatabaseConfig struct {
	Host            string `koanf:"host" validate:"required"`
	Port            int    `koanf:"port" validate:"required"`
	User            string `koanf:"user" validate:"required"`
	Password        string `koanf:"password" validate:"required"`
	Name            string `koanf:"name" validate:"required"`
	SSLMode         string `koanf:"ssl_mode" validate:"required"`
	MaxOpenConns    int    `koanf:"max_open_conns" validate:"required"`
	MaxIdleConns    int    `koanf:"max_idle_conns" validate:"required"`
	ConnMaxLifetime int    `koanf:"conn_max_lifetime" validate:"required"`
	ConnMaxIdleTime int    `koanf:"conn_max_idle_time" validate:"required"`
}

type AuthConfig struct {
	SecretKey     string `koanf:"secret_key" validate:"required"`
	WebhookSecret string `koanf:"webhook_secret" validate:"required"`
}

type RedisConfig struct {
	Address string `koanf:"address" validate:"required"`
}

type IntegrationConfig struct {
	SMTPHost     string `koanf:"smtp_host" validate:"required"`
	SMTPPort     string `koanf:"smtp_port" validate:"required"`
	SMTPUsername string `koanf:"smtp_username" validate:"required"`
	SMTPPassword string `koanf:"smtp_password" validate:"required"`
	SMTPFrom     string `koanf:"smtp_from" validate:"required"`
}

type AWSConfig struct {
	Region          string `koanf:"region"`
	AccessKeyID     string `koanf:"access_key_id"`
	SecretAccessKey string `koanf:"secret_access_key"`
	EndpointURL     string `koanf:"endpoint_url"`
	UploadBucket    string `koanf:"upload_bucket"`
}

func (c AWSConfig) IsConfigured() bool {
	return c.AccessKeyID != "" && c.SecretAccessKey != "" && c.UploadBucket != ""
}

type SearchConfig struct {
	URL      string `koanf:"url"`
	Username string `koanf:"username"`
	Password string `koanf:"password"`
}

func (c SearchConfig) IsConfigured() bool {
	return c.URL != ""
}

func LoadConfig() (*Config, error) {
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()

	k := koanf.New(".")

	err := k.Load(env.Provider("NEXUS_", ".", func(s string) string {
		return strings.ToLower(strings.TrimPrefix(s, "NEXUS_"))
	}), nil)

	if err != nil {
		logger.Fatal().Err(err).Msg("could not load initial env variables")
	}

	mainConfig := &Config{}

	err = k.Unmarshal("", mainConfig)

	if platformPort := os.Getenv("PORT"); platformPort != "" {
		mainConfig.Server.Port = platformPort
	}

	validate := validator.New()

	err = validate.Struct(mainConfig)

	if err != nil {
		logger.Fatal().Err(err).Msg("config validation failed")
	}

	return mainConfig, nil

}
