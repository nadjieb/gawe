package gawe

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"
)

//go:generate mockery --name=IdentifiableJob --structname=IdentifiableJob --filename=identifiable_job.go --output=gawetest --outpkg=gawetest

// IdentifiableJob defines the interface a job must have to be identifiable
type IdentifiableJob interface {
	JobID() string
	JobType() string
	Tags() []string
}

//go:generate mockery --name=Job --structname=Job --filename=job.go --output=gawetest --outpkg=gawetest

// Job defines the interface a job must have
type Job interface {
	JobID() string
	JobType() string
	Tags() []string
	Exec(ctx context.Context) error
}

//go:generate mockery --name=Engineable --structname=Engine --filename=engine.go --output=gawetest --outpkg=gawetest

// Engineable defines the interface a workers controller must have
type Engineable interface {
	Start()
	Stop()
	Enqueue(ctx context.Context, job Job) error
}

// Engine is the workers controller
type Engine struct {
	cfn               chan func() (context.Context, Job)
	cw                chan struct{}
	wg                sync.WaitGroup
	maxAttempts       int
	maxQueueSize      int
	maxWorkers        int
	inactivityTimeout time.Duration
	plugins           []Plugin
}

var _ Engineable = (*Engine)(nil)

const (
	defaultMaxAttempts       = 1
	defaultMaxQueueSize      = 100
	defaultMaxWorkers        = 1
	defaultInactivityTimeout = 30 * time.Second
)

// NewEngine returns a new instance of workers controller
func NewEngine(opts ...Option) *Engine {
	e := &Engine{
		maxAttempts:       defaultMaxAttempts,
		maxQueueSize:      defaultMaxQueueSize,
		maxWorkers:        defaultMaxWorkers,
		inactivityTimeout: defaultInactivityTimeout,
	}

	for _, opt := range opts {
		opt(e)
	}

	return e
}

// Start opening channels in the workers controller
func (e *Engine) Start() {
	e.cfn = make(chan func() (context.Context, Job), e.maxQueueSize)
	e.cw = make(chan struct{}, e.maxWorkers)
}

// Stop closing channels in the workers controller
func (e *Engine) Stop() {
	close(e.cfn)
	e.wg.Wait()
	close(e.cw)
}

// Enqueue put the job in the queue unless it is full
func (e *Engine) Enqueue(ctx context.Context, job Job) error {
	// Run a new worker unless it is full
	select {
	case e.cw <- struct{}{}:
		e.runNewWorker()
	default:
	}

	// Put the payload in the channel unless it is full
	select {
	case e.cfn <- func() (context.Context, Job) { return ctx, job }:
		return nil
	default:
		return errors.New("The job's queue is full")
	}
}

func (e *Engine) executeJob(ctx context.Context, job Job) {
	currentCtx := ctx

	for i := 0; i < e.maxAttempts; i++ {
		currentCtx = e.reportJobStart(currentCtx, job)

		if err := job.Exec(currentCtx); err != nil {
			currentCtx = e.reportJobError(currentCtx, job, err)
			continue
		}

		e.reportJobEnd(currentCtx, job)
		break
	}
}

func (e *Engine) runNewWorker() {
	e.wg.Add(1)

	go func() {
		defer func() {
			<-e.cw
			e.wg.Done()

			if err := recover(); err != nil {
				log.Println(err)
			}
		}()

		for {
			select {
			case <-time.After(e.inactivityTimeout):
				return
			case fn := <-e.cfn:
				if fn == nil {
					return
				}
				e.executeJob(fn())
			}
		}
	}()
}

func (e *Engine) reportJobStart(ctx context.Context, job IdentifiableJob) context.Context {
	currentCtx := ctx

	for _, p := range e.plugins {
		currentCtx = p.OnJobStart(currentCtx, job)
	}

	return currentCtx
}

func (e *Engine) reportJobEnd(ctx context.Context, job IdentifiableJob) {
	for _, p := range e.plugins {
		p.OnJobEnd(ctx, job)
	}
}

func (e *Engine) reportJobError(ctx context.Context, job IdentifiableJob, err error) context.Context {
	currentCtx := ctx

	for _, p := range e.plugins {
		currentCtx = p.OnJobError(currentCtx, job, err)
	}

	return currentCtx
}
