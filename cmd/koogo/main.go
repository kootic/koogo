package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"github.com/kootic/koogo/internal/app"
	"github.com/kootic/koogo/internal/config"
	"github.com/kootic/koogo/internal/jobs"
	_ "github.com/kootic/koogo/swagger"
)

//	@title						Kootic Starter Project
//	@version					0.0.1
//	@description				This is a boilerplate for Go API projects.
//	@contact.name				Alex
//	@contact.url				https://github.com/kootic/koogo
//	@contact.email				alex@kootic.com
//	@license.name				Apache 2.0
//	@license.url				http://www.apache.org/licenses/LICENSE-2.0.html
//	@host						<host>
//	@BasePath					/api
//	@securityDefinitions.basic	BasicAuth.

var rootCmd = &cobra.Command{
	Use:   "koogo",
	Short: "Koogo is a production-ready Go API",
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the server",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Load config
		cfg, err := config.LoadConfigFromEnv("", true)
		if err != nil {
			return err
		}

		app := app.NewApp(cfg)

		// Setup signal handling for graceful shutdown
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		// Bootstrap application
		if err := app.Bootstrap(ctx); err != nil {
			return err
		}

		// Start application (blocks until shutdown signal)
		if err := app.Start(ctx); err != nil {
			return err
		}

		// Create a timeout context for shutdown operations
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Shutdown application gracefully
		return app.Shutdown(shutdownCtx)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	// Dynamically create job commands from JobsRegistry
	for jobID, job := range jobs.JobsRegistry {
		jobCmd := &cobra.Command{
			Use:   jobID,
			Short: "Run " + jobID + " job",
			RunE: func(cmd *cobra.Command, args []string) error {
				// Load config
				cfg, err := config.LoadConfigFromEnv("", false)
				if err != nil {
					return err
				}

				// Setup signal handling
				ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
				defer stop()

				// Collect all flag values
				flags := make(map[string]string)
				for _, requiredFlag := range job.RequiredFlags {
					value, err := cmd.Flags().GetString(requiredFlag)
					if err != nil {
						return err
					}

					flags[requiredFlag] = value
				}

				for _, optionalFlag := range job.OptionalFlags {
					value, err := cmd.Flags().GetString(optionalFlag)
					if err != nil {
						return err
					}

					if value != "" {
						flags[optionalFlag] = value
					}
				}

				// Run the job
				return jobs.RunJob(ctx, cfg, jobID, flags)
			},
		}

		// Add required flags
		for _, flag := range job.RequiredFlags {
			jobCmd.Flags().String(flag, "", "Required flag for "+jobID+" job")

			err := jobCmd.MarkFlagRequired(flag)
			if err != nil {
				log.Fatal(err)
			}
		}

		// Add optional flags
		for _, flag := range job.OptionalFlags {
			jobCmd.Flags().String(flag, "", "Optional flag for "+jobID+" job")
		}

		// Add the job command to the root command
		rootCmd.AddCommand(jobCmd)
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
