package app

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/contrib/bridges/otelzap"
	"go.opentelemetry.io/otel/log/global"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/kootic/koogo/internal/config"
	"github.com/kootic/koogo/internal/server"
	"github.com/kootic/koogo/pkg/koodb"
	"github.com/kootic/koogo/pkg/koolog"
	"github.com/kootic/koogo/pkg/kootel"
)

// App represents the application and its dependencies.
type App struct {
	config       *config.Config
	logger       *zap.Logger
	fiberApp     *fiber.App
	server       server.Server
	cleanupFuncs []func(ctx context.Context) error
}

// NewApp creates a new App instance.
func NewApp(cfg *config.Config) *App {
	return &App{
		config: cfg,
	}
}

// Bootstrap initializes the application and its dependencies.
func (a *App) Bootstrap(ctx context.Context) error {
	// Initialize logger
	logger, err := koolog.NewLogger(a.config.App.IsProd(), a.config.App.ZapLogLevel())
	if err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}

	// Initialize OpenTelemetry
	stop, err := kootel.InitializeOTel(ctx, kootel.OTelConfig{
		ServiceName:    a.config.App.Name,
		ServiceVersion: a.config.App.Version,
		Environment:    string(a.config.App.Env),
		ExporterType:   a.config.OTel.Exporter,
	})
	if err != nil {
		return fmt.Errorf("failed to initialize OpenTelemetry: %w", err)
	}

	a.cleanupFuncs = append(a.cleanupFuncs, stop)

	a.logger = zap.New(
		zapcore.NewTee(
			logger.Core(),
			otelzap.NewCore(a.config.App.Name, otelzap.WithLoggerProvider(global.GetLoggerProvider())),
		),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)

	// Initialize database client
	var sqldb *sql.DB
	if a.config.App.IsTest() {
		sqldb, err = koodb.NewPostgresTxDB(a.config.Database.DSN())
	} else {
		poolConfig := &koodb.PoolConfig{
			MaxConns:          a.config.Database.MaxConns,
			MinConns:          a.config.Database.MinConns,
			MaxConnLifetime:   time.Duration(a.config.Database.MaxConnLifetime) * time.Minute,
			MaxConnIdleTime:   time.Duration(a.config.Database.MaxConnIdleTime) * time.Minute,
			ConnectionTimeout: time.Duration(a.config.Database.ConnectionTimeout) * time.Second,
		}
		sqldb, err = koodb.NewPostgresPool(ctx, a.config.Database.DSN(), poolConfig)
	}

	if err != nil {
		return fmt.Errorf("failed to create database pool: %w", err)
	}

	a.cleanupFuncs = append(a.cleanupFuncs, func(ctx context.Context) error {
		return sqldb.Close()
	})

	// Create server with fiber app
	a.fiberApp = fiber.New(fiber.Config{
		ReadTimeout:  time.Duration(a.config.App.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(a.config.App.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(a.config.App.IdleTimeout) * time.Second,
		BodyLimit:    a.config.App.BodyLimit * 1024 * 1024, // Convert MB to bytes
	})

	srv, err := server.NewServer(a.config, a.logger, sqldb, a.fiberApp)
	if err != nil {
		return fmt.Errorf("failed to create server: %w", err)
	}

	a.server = srv

	// Initialize fiber app
	err = a.server.Initialize()
	if err != nil {
		return fmt.Errorf("failed to initialize fiber app: %w", err)
	}

	return nil
}

// Start starts the application and blocks until shutdown signal is received.
func (a *App) Start(ctx context.Context) error {
	// Start server in a goroutine
	serverErr := make(chan error, 1)
	go func() {
		if err := a.server.Start(); err != nil {
			serverErr <- fmt.Errorf("failed to start server: %w", err)
		}
	}()

	a.logger.Info("Application started successfully")

	// Wait for shutdown signal or server error
	select {
	case <-ctx.Done():
		a.logger.Info("Shutdown signal received, gracefully shutting down...")
		return nil
	case err := <-serverErr:
		return err
	}
}

// Shutdown gracefully shuts down the application.
func (a *App) Shutdown(ctx context.Context) error {
	a.logger.Info("Initiating graceful shutdown...")

	if a.server != nil {
		if err := a.server.Shutdown(ctx); err != nil {
			a.logger.Error("Failed to shutdown server", zap.Error(err))
			return fmt.Errorf("failed to shutdown server: %w", err)
		}

		a.logger.Info("Server shutdown complete")
	}

	for i, fn := range a.cleanupFuncs {
		if err := fn(ctx); err != nil {
			a.logger.Error("Failed to cleanup app", zap.Error(err), zap.Int("cleanup_index", i))
			return fmt.Errorf("failed to cleanup: %w", err)
		}
	}

	a.logger.Info("Application shutdown successfully")

	return nil
}

func (a *App) FiberApp() *fiber.App {
	if a.fiberApp == nil {
		return nil
	}

	return a.fiberApp
}
