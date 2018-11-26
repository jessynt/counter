package counter

import (
	"context"
)

type Counter interface {
	Incr(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (int64, error)
}
