package filter

import (
	"context"
)

type KV interface {
	Get(ctx context.Context) ([]byte, error)
	Set(ctx context.Context, value []byte) error
}
