package repo

import "context"

type HealthRepository interface {
	Ping(ctx context.Context) error
}
