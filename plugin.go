package gawe

import (
	"context"
)

// Plugin defines the interface that a plugin must have
type Plugin interface {
	OnJobStart(ctx context.Context, job IdentifiableJob) context.Context
	OnJobEnd(ctx context.Context, job IdentifiableJob)
	OnJobError(ctx context.Context, job IdentifiableJob, err error) context.Context
}
