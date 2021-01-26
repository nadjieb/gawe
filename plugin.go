package gawe

import (
	"context"
)

//go:generate mockery --name=Plugin --structname=Plugin --filename=plugin.go --output=gawetest --outpkg=gawetest

// Plugin defines the interface that a plugin must have
type Plugin interface {
	OnJobStart(ctx context.Context, job IdentifiableJob) context.Context
	OnJobEnd(ctx context.Context, job IdentifiableJob)
	OnJobError(ctx context.Context, job IdentifiableJob, err error) context.Context
}
