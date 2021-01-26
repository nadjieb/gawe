package gawe

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type pluginExampleCtxKey uint8

const (
	pluginExampleKey pluginExampleCtxKey = iota
	pluginExample2Key
)

type pluginExample struct{}

var _ Plugin = (*pluginExample)(nil)

func (p *pluginExample) OnJobStart(ctx context.Context, job IdentifiableJob) context.Context {
	if valCtx := ctx.Value(pluginExampleKey); valCtx != nil {
		if val := valCtx.(string); val != "fromError" {
			panic(errors.New("OnJobStart"))
		}
	}

	if valCtx := ctx.Value(pluginExample2Key); valCtx == nil {
		panic(errors.New("OnJobStart"))
	} else if val := valCtx.(string); val != "beforeStart" {
		panic(errors.New("OnJobStart"))
	}

	return context.WithValue(ctx, pluginExampleKey, "fromStart")
}

func (p *pluginExample) OnJobEnd(ctx context.Context, job IdentifiableJob) {
	if valCtx := ctx.Value(pluginExampleKey); valCtx == nil {
		panic(errors.New("OnJobEnd"))
	} else if val := valCtx.(string); val != "fromStart" {
		panic(errors.New("OnJobError"))
	}
}

func (p *pluginExample) OnJobError(ctx context.Context, job IdentifiableJob, err error) context.Context {
	if valCtx := ctx.Value(pluginExampleKey); valCtx == nil {
		panic(errors.New("OnJobError"))
	} else if val := valCtx.(string); val != "fromStart" {
		panic(errors.New("OnJobError"))
	}

	return context.WithValue(ctx, pluginExampleKey, "fromError")
}

func TestPlugin(t *testing.T) {
	t.Parallel()

	engine := NewEngine(
		WithMaxAttempts(2),
		WithPlugins(&pluginExample{}),
	)
	assert.Equal(t, 2, engine.maxAttempts)
	assert.Equal(t, 1, len(engine.plugins))

	engine.Start()
	defer engine.Stop()

	ctx := context.Background()
	ctx = context.WithValue(ctx, pluginExample2Key, "beforeStart")

	err := engine.Enqueue(ctx, &jobExample{})
	time.Sleep(100 * time.Millisecond)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(engine.cw))

	err = engine.Enqueue(ctx, &jobExample{IsError: true})
	time.Sleep(100 * time.Millisecond)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(engine.cw))
}
