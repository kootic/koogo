package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/kootic/koogo/internal/app"
	"github.com/kootic/koogo/internal/config"
	_ "github.com/kootic/koogo/swagger"
)

//go:generate swag fmt
//go:generate swag init --parseDependency --parseInternal --output ../../swagger

//	@title			Kootic Starter Project
//	@version		0.0.1
//	@description	This is a boilerplate for Go API projects.

//	@contact.name	Alex
//	@contact.url	https://github.com/kootic/koogo
//	@contact.email	alex@kootic.com

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		<host>
//	@BasePath	/api

// @securityDefinitions.basic	BasicAuth.
func main() {
	// Load config
	cfg, err := config.LoadConfigFromEnv("")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	app := app.NewApp(cfg)

	// Setup signal handling
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Bootstrap application
	if err := app.Bootstrap(ctx); err != nil {
		log.Fatalf("Failed to bootstrap application: %v", err) //nolint:gocritic
	}

	// Start application
	if err := app.Start(ctx); err != nil {
		log.Fatalf("Application failed: %v", err)
	}

	// Shutdown application
	if err := app.Shutdown(ctx); err != nil {
		log.Fatalf("Failed to shutdown application: %v", err)
	}
}
