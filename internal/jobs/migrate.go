package jobs

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"ariga.io/atlas-go-sdk/atlasexec"

	"github.com/kootic/koogo/internal/config"
)

func Migrate(ctx context.Context, cfg *config.Config, flags map[string]string) error {
	migrationsDir, ok := flags["migrations-dir"]
	if !ok {
		return fmt.Errorf("--migrations-dir is required")
	}

	if err := applyMigrations(ctx, migrationsDir, cfg.Database.DSN()); err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	return nil
}

func applyMigrations(ctx context.Context, migrationsDir string, dbURL string) error {
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
		URL: dbURL,
	})
	if err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	for _, migration := range result.Applied {
		log.Printf("Applied migration: %s", migration.Name)
	}

	return nil
}
