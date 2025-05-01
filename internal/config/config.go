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
	Name         string
	Version      string
	Env          AppEnv
	Port         int
	LogLevel     AppLogLevel
	ReadTimeout  int // Read timeout in seconds
	WriteTimeout int // Write timeout in seconds
	IdleTimeout  int // Idle timeout in seconds
	BodyLimit    int // Body limit in megabytes
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
	Enabled  bool
	Exporter kootel.OTelExporterType
}

func (o *OTelConfig) Validate() error {
	if o.Enabled && o.Exporter == "" {
		return fmt.Errorf("otel config is incomplete")
	}

	return nil
}

type DatabaseConfig struct {
	Host              string
	Port              int
	Username          string
	Password          string
	Database          string
	MaxConns          int    // Maximum number of connections in the pool
	MinConns          int    // Minimum number of connections in the pool
	MaxConnLifetime   int    // Maximum lifetime of a connection in minutes
	MaxConnIdleTime   int    // Maximum idle time of a connection in minutes
	ConnectionTimeout int    // Connection timeout in seconds
	SSLMode           string // SSL mode for the database connection
}

func (d *DatabaseConfig) Validate() error {
	if d.Host == "" || d.Port == 0 || d.Username == "" || d.Password == "" {
		return fmt.Errorf("database config is incomplete")
	}

	return nil
}

func (d *DatabaseConfig) DSN() string {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", d.Username, d.Password, d.Host, d.Port, d.Database)
	if d.SSLMode != "" {
		dsn += "?sslmode=" + d.SSLMode
	}

	return dsn
}

func (d *DatabaseConfig) DSNWithoutDatabase() string {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d", d.Username, d.Password, d.Host, d.Port)
	if d.SSLMode != "" {
		dsn += "?sslmode=" + d.SSLMode
	}

	return dsn
}

// getEnvAsInt parses an environment variable as int with a default value.
func getEnvAsInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return intValue
}

// getEnvAsString parses an environment variable as string with a default value.
func getEnvAsString(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	return value
}

// getEnvAsBool parses an environment variable as boolean with a default value.
func getEnvAsBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}

	return boolValue
}

// LoadConfigFromEnv prepares the config from environment variables. When not running locally,
// we must set KOO_APP_ENV to a valid environment. If KOO_APP_ENV is not set or is set to "local",
// we will load the .env file at envFilePath or simply use the .env file in the same directory as the main.go file.
func LoadConfigFromEnv(envFilePath string, validate bool) (*Config, error) {
	appEnv := AppEnv(getEnvAsString("KOO_APP_ENV", string(AppEnvLocal)))

	// This is only for local development, env files will not be included in the build and we rely on environment variables
	if appEnv == AppEnvLocal {
		if envFilePath == "" {
			envFilePath = ".env"
		}

		if err := godotenv.Load(envFilePath); err != nil {
			log.Println("Unable to load .env file, using environment variables only; error:", err)
		}
	}

	var err error

	var appPort int

	appPort, err = strconv.Atoi(os.Getenv("KOO_APP_PORT"))
	if validate && err != nil {
		return nil, fmt.Errorf("invalid app port: %w", err)
	}

	appConfig := AppConfig{
		Name:         os.Getenv("KOO_APP_NAME"),
		Version:      os.Getenv("KOO_APP_VERSION"),
		Env:          appEnv,
		Port:         appPort,
		LogLevel:     AppLogLevel(os.Getenv("KOO_APP_LOG_LEVEL")),
		ReadTimeout:  getEnvAsInt("KOO_APP_READ_TIMEOUT_SECONDS", 15),
		WriteTimeout: getEnvAsInt("KOO_APP_WRITE_TIMEOUT_SECONDS", 15),
		IdleTimeout:  getEnvAsInt("KOO_APP_IDLE_TIMEOUT_SECONDS", 120),
		BodyLimit:    getEnvAsInt("KOO_APP_BODY_LIMIT_MB", 4),
	}

	swaggerConfig := SwaggerConfig{
		Enabled:  getEnvAsBool("KOO_SWAGGER_ENABLED", false),
		Username: os.Getenv("KOO_SWAGGER_USERNAME"),
		Password: os.Getenv("KOO_SWAGGER_PASSWORD"),
	}

	oTelConfig := OTelConfig{
		Enabled:  os.Getenv("KOO_OTEL_ENABLED") == "true",
		Exporter: kootel.OTelExporterType(os.Getenv("KOO_OTEL_EXPORTER")),
	}

	var dbPort int

	dbPort, err = strconv.Atoi(os.Getenv("KOO_DB_PORT"))
	if validate && err != nil {
		return nil, fmt.Errorf("invalid db port: %w", err)
	}

	// Parse connection pool settings with defaults
	maxConns := getEnvAsInt("KOO_DB_MAX_CONNS", 25)
	minConns := getEnvAsInt("KOO_DB_MIN_CONNS", 5)
	maxConnLifetime := getEnvAsInt("KOO_DB_MAX_CONN_LIFETIME_MINUTES", 60)
	maxConnIdleTime := getEnvAsInt("KOO_DB_MAX_CONN_IDLE_TIME_MINUTES", 30)
	connectionTimeout := getEnvAsInt("KOO_DB_CONNECTION_TIMEOUT_SECONDS", 10)

	databaseConfig := DatabaseConfig{
		Host:              os.Getenv("KOO_DB_HOST"),
		Port:              dbPort,
		Username:          os.Getenv("KOO_DB_USERNAME"),
		Password:          os.Getenv("KOO_DB_PASSWORD"),
		Database:          os.Getenv("KOO_DB_DATABASE"),
		MaxConns:          maxConns,
		MinConns:          minConns,
		MaxConnLifetime:   maxConnLifetime,
		MaxConnIdleTime:   maxConnIdleTime,
		ConnectionTimeout: connectionTimeout,
	}

	config := Config{
		App:      appConfig,
		Swagger:  swaggerConfig,
		OTel:     oTelConfig,
		Database: databaseConfig,
	}

	if validate {
		err = config.Validate()
		if err != nil {
			return nil, fmt.Errorf("invalid config: %w", err)
		}
	}

	return &config, nil
}
