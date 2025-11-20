package communication

import (
	"context"
	"time"
)

type CheckFunc func(ctx context.Context, t Target) (bool, error)

type Target struct {
	Host    string
	Port    int
	Path    string
	Timeout time.Duration
	Service string // for grpc
}
