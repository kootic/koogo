package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"go.uber.org/zap/zapcore"

	"github.com/kootic/koogo/pkg/kootel"
)

type AppEnv string

const (
	AppEnvTest    AppEnv = "test"
	AppEnvLocal   AppEnv = "local"
	AppEnvDev     AppEnv = "dev"
	AppEnvStaging AppEnv = "staging"
	AppEnvProd    AppEnv = "prod"
)

var validEnvs = map[AppEnv]bool{
	AppEnvTest:    true,
	AppEnvLocal:   true,
	AppEnvDev:     true,
	AppEnvStaging: true,
	AppEnvProd:    true,
}

type AppLogLevel string

const (
	AppLogLevelDebug AppLogLevel = "debug"
	AppLogLevelInfo  AppLogLevel = "info"
	AppLogLevelWarn  AppLogLevel = "warn"
	AppLogLevelError AppLogLevel = "error"
)

type Config struct {
	App      AppConfig
	Swagger  SwaggerConfig
	OTel     OTelConfig
	Database DatabaseConfig
}

func (c *Config) Validate() error {
	if err := c.App.Validate(); err != nil {
		return fmt.Errorf("app config is invalid: %w", err)
	}

	if err := c.Swagger.Validate(); err != nil {
		return fmt.Errorf("swagger config is invalid: %w", err)
	}

	if err := c.OTel.Validate(); err != nil {
		return fmt.Errorf("otel config is invalid: %w", err)
	}

	if err := c.Database.Validate(); err != nil {
		return fmt.Errorf("database config is invalid: %w", err)
	}

	return nil
}

type AppConfig struct {
	Name     string
	Version  string
	Env      AppEnv
	Port     int
	LogLevel AppLogLevel
}

func (a *AppConfig) Validate() error {
	if a.Name == "" || a.Version == "" || a.Env == "" || a.Port == 0 || a.LogLevel == "" {
		return fmt.Errorf("app config is incomplete")
	}

	if !validEnvs[a.Env] {
		return fmt.Errorf("invalid app env: %s", a.Env)
	}

	return nil
}

func (a *AppConfig) IsProd() bool {
	return a.Env == AppEnvProd
}

func (a *AppConfig) IsTest() bool {
	return a.Env == AppEnvTest
}

func (a *AppConfig) ZapLogLevel() zapcore.Level {
	switch a.LogLevel {
	case AppLogLevelDebug:
		return zapcore.DebugLevel
	case AppLogLevelInfo:
		return zapcore.InfoLevel
	case AppLogLevelWarn:
		return zapcore.WarnLevel
	case AppLogLevelError:
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

type SwaggerConfig struct {
	Enabled  bool
	Username string
	Password string
}

func (s *SwaggerConfig) Validate() error {
	if s.Enabled && (s.Username == "" || s.Password == "") {
		return fmt.Errorf("swagger config is incomplete")
	}

	return nil
}

type OTelConfig struct {
	Enabled      bool
	Exporter     kootel.OTelExporterType
	OTLPEndpoint string
}

func (o *OTelConfig) Validate() error {
	if o.Enabled && (o.Exporter == "" || o.OTLPEndpoint == "") {
		return fmt.Errorf("otel config is incomplete")
	}

	return nil
}

type DatabaseConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
}

func (d *DatabaseConfig) Validate() error {
	if d.Host == "" || d.Port == 0 || d.Username == "" || d.Password == "" {
		return fmt.Errorf("database config is incomplete")
	}

	return nil
}

func (d *DatabaseConfig) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", d.Username, d.Password, d.Host, d.Port, d.Database)
}

func (d *DatabaseConfig) DSNWithoutDatabase() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d?sslmode=disable", d.Username, d.Password, d.Host, d.Port)
}

// LoadConfigFromEnv prepares the config from environment variables. When not running locally,
// we must set KOO_APP_ENV to a valid environment. If KOO_APP_ENV is not set or is set to "local",
// we will load the .env file at envFilePath or simply use the .env file in the same directory as the main.go file.
func LoadConfigFromEnv(envFilePath string) (*Config, error) {
	appEnv := AppEnv(os.Getenv("KOO_APP_ENV"))

	// This is only for local development, env files will not be included in the build and we rely on environment variables
	if appEnv == "" || appEnv == AppEnvLocal {
		appEnv = AppEnvLocal

		if envFilePath == "" {
			envFilePath = ".env"
		}

		if err := godotenv.Load(envFilePath); err != nil {
			log.Println("Unable to load .env file, using environment variables only")
		}
	}

	appPort, err := strconv.Atoi(os.Getenv("KOO_APP_PORT"))
	if err != nil {
		return nil, fmt.Errorf("invalid app port: %w", err)
	}

	appConfig := AppConfig{
		Name:     os.Getenv("KOO_APP_NAME"),
		Version:  os.Getenv("KOO_APP_VERSION"),
		Env:      appEnv,
		Port:     appPort,
		LogLevel: AppLogLevel(os.Getenv("KOO_APP_LOG_LEVEL")),
	}

	swaggerConfig := SwaggerConfig{
		Enabled:  os.Getenv("KOO_SWAGGER_ENABLED") == "true",
		Username: os.Getenv("KOO_SWAGGER_USERNAME"),
		Password: os.Getenv("KOO_SWAGGER_PASSWORD"),
	}

	oTelConfig := OTelConfig{
		Enabled:      os.Getenv("KOO_OTEL_ENABLED") == "true",
		Exporter:     kootel.OTelExporterType(os.Getenv("KOO_OTEL_EXPORTER")),
		OTLPEndpoint: os.Getenv("KOO_OTEL_OTLP_ENDPOINT"),
	}

	dbPort, err := strconv.Atoi(os.Getenv("KOO_DB_PORT"))
	if err != nil {
		return nil, fmt.Errorf("invalid db port: %w", err)
	}

	databaseConfig := DatabaseConfig{
		Host:     os.Getenv("KOO_DB_HOST"),
		Port:     dbPort,
		Username: os.Getenv("KOO_DB_USERNAME"),
		Password: os.Getenv("KOO_DB_PASSWORD"),
		Database: os.Getenv("KOO_DB_DATABASE"),
	}

	config := Config{
		App:      appConfig,
		Swagger:  swaggerConfig,
		OTel:     oTelConfig,
		Database: databaseConfig,
	}

	err = config.Validate()
	if err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &config, nil
}
