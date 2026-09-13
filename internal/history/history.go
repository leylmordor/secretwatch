package history

import (
	"context"
	"time"
)

type Entry struct {
	Timestamp  time.Time
	Actor      string
	Action     string
	Source     string // cloudtrail, gcp-audit, github-audit
	SecretName string
	Region     string
	Note       string // extra context, e.g. "automated rotation" vs "manual"
}

type Fetcher interface {
	StoreType() string
	Fetch(ctx context.Context, secretName string, limit int) ([]*Entry, error)
}
