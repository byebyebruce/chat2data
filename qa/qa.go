package qa

import (
	"context"
)

type QA interface {
	Answer(ctx context.Context, question string) (string, error)
}
