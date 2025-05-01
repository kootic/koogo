// Package tests contains integration tests for the application.
// Integration tests are only enabled when the RUN_INTEGRATION_TESTS environment variable is set to true.
package tests

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"ariga.io/atlas-go-sdk/atlasexec"

	"github.com/kootic/koogo/internal/app"
	"github.com/kootic/koogo/internal/config"
	"github.com/kootic/koogo/pkg/koodb"
)

const (
	migrationsDir         = "../repository/dbrepo/migrations"
	initializationTimeout = 10 * time.Second
	testTimeout           = 60 * time.Second
)

var (
	testConfig *config.Config
	testApp    *app.App
)

// TestMain is the entry point for the integration tests.
func TestMain(m *testing.M) {
	os.Exit(runIntegrationTests(m))
}

// runIntegrationTests runs the integration tests and returns the exit code
// if the RUN_INTEGRATION_TESTS environment variable is set to true. If not,
// it returns 0 and skips the tests.
func runIntegrationTests(m *testing.M) int {
	enabled := strings.ToLower(os.Getenv("RUN_INTEGRATION_TESTS")) == "true"
	if !enabled {
		return 0
	}

	var err error

	testConfig, err = config.LoadConfigFromEnv(".env.test")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	dbCleanup, err := initializeTestDB(testConfig)
	if err != nil {
		log.Fatalf("Failed to initialize test database: %v", err) //nolint:gocritic
	}
	defer dbCleanup(ctx) //nolint:errcheck

	testApp = app.NewApp(testConfig)
	defer testApp.Shutdown(ctx) //nolint:errcheck

	err = testApp.Bootstrap(ctx)
	if err != nil {
		log.Fatalf("Failed to bootstrap app: %v", err)
	}

	return m.Run()
}

func initializeTestDB(config *config.Config) (func(ctx context.Context) error, error) {
	// Generate a unique database name for the test
	config.Database.Database = config.App.Name + "test" + strconv.FormatInt(time.Now().Unix(), 10)

	ctx, cancel := context.WithTimeout(context.Background(), initializationTimeout)
	defer cancel()

	sqlDB, err := koodb.NewPostgresConn(ctx, config.Database.DSNWithoutDatabase())
	if err != nil {
		return nil, err
	}
	defer sqlDB.Close() //nolint:errcheck

	_, err = sqlDB.ExecContext(ctx, fmt.Sprintf("DROP DATABASE IF EXISTS %s", config.Database.Database))
	if err != nil {
		return nil, err
	}

	_, err = sqlDB.ExecContext(ctx, fmt.Sprintf("CREATE DATABASE %s", config.Database.Database))
	if err != nil {
		return nil, err
	}

	err = applyMigrations(ctx, config.Database.DSN())
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) error {
		cleanupTestDB(ctx, config)
		return nil
	}, nil
}

func applyMigrations(ctx context.Context, testDBURL string) error {
	absMigrationsDir, err := filepath.Abs(migrationsDir)
	if err != nil {
		return err
	}

	workdir, err := atlasexec.NewWorkingDir(
		atlasexec.WithMigrations(
			os.DirFS(absMigrationsDir),
		),
	)
	if err != nil {
		return fmt.Errorf("failed to load working directory: %w", err)
	}

	defer func() {
		_ = workdir.Close()
	}()

	client, err := atlasexec.NewClient(workdir.Path(), "atlas")
	if err != nil {
		return fmt.Errorf("failed to initialize client: %w", err)
	}

	result, err := client.MigrateApply(ctx, &atlasexec.MigrateApplyParams{
		URL: testDBURL,
	})
	if err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	for _, migration := range result.Applied {
		log.Printf("Applied migration: %s", migration.Name)
	}

	return nil
}

func cleanupTestDB(ctx context.Context, config *config.Config) {
	sqlDB, err := koodb.NewPostgresConn(ctx, config.Database.DSNWithoutDatabase())
	if err != nil {
		log.Printf("Failed to connect to test database: %v", err)
	}
	defer sqlDB.Close() //nolint:errcheck

	_, err = sqlDB.ExecContext(ctx, fmt.Sprintf("DROP DATABASE IF EXISTS %s", config.Database.Database))
	if err != nil {
		log.Printf("Failed to drop test database: %v", err)
	}
}
