package jobs

import (
	"context"
	"fmt"

	"github.com/kootic/koogo/internal/config"
	"github.com/kootic/koogo/pkg/kooctx"
	"github.com/kootic/koogo/pkg/koolog"
)

type Job struct {
	RequiredFlags []string
	OptionalFlags []string
	Run           func(ctx context.Context, cfg *config.Config, flags map[string]string) error
}

var JobsRegistry = map[string]Job{
	"migrate": {
		RequiredFlags: []string{"migrations-dir"},
		Run:           Migrate,
	},
}

func RunJob(ctx context.Context, cfg *config.Config, jobID string, flags map[string]string) error {
	logger, err := koolog.NewLogger(cfg.App.IsProd(), cfg.App.ZapLogLevel())
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}

	kooctx.SetContextLogger(ctx, logger)

	job, ok := JobsRegistry[jobID]
	if !ok {
		return fmt.Errorf("job %s not found", jobID)
	}

	return job.Run(ctx, cfg, flags)
}
