package containerization

import (
	"context"
	"io"
)

type Client interface {
	GetLogs(ctx context.Context) (io.ReadCloser, error)
}
