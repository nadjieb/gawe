package gawe

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type jobExample struct {
	ID            string
	IsError       bool
	IsPanic       bool
	SleepDuration *time.Duration
	AttemptsCount int
}

var _ Job = (*jobExample)(nil)

func (j *jobExample) JobID() string {
	return j.ID
}

func (j *jobExample) JobType() string {
	return "job-example"
}

func (j *jobExample) Tags() []string {
	return []string{"job", "example"}
}

func (j *jobExample) Exec(ctx context.Context) error {
	j.AttemptsCount++

	if j.IsPanic {
		panic(errors.New("panic"))
	}

	if j.IsError {
		return errors.New("error")
	}

	if j.SleepDuration != nil {
		time.Sleep(*j.SleepDuration)
	}

	return nil
}

func TestEngine(t *testing.T) {
	t.Parallel()

	engine := NewEngine()
	assert.Equal(t, defaultMaxAttempts, engine.maxAttempts)
	assert.Equal(t, defaultMaxQueueSize, engine.maxQueueSize)
	assert.Equal(t, defaultMaxWorkers, engine.maxWorkers)
	assert.Equal(t, defaultInactivityTimeout, engine.inactivityTimeout)
	assert.Nil(t, engine.plugins)

	engine.Start()

	err := engine.Enqueue(context.TODO(), &jobExample{})
	assert.Nil(t, err)
	assert.Equal(t, 1, len(engine.cw))

	engine.Stop()
	assert.Equal(t, 0, len(engine.cw))
	assert.Panics(t, func() { close(engine.cfn) })
	assert.Panics(t, func() { close(engine.cw) })
}

func TestPanic(t *testing.T) {
	t.Parallel()

	engine := NewEngine()

	engine.Start()
	defer engine.Stop()

	err := engine.Enqueue(context.TODO(), &jobExample{IsPanic: true})
	assert.Nil(t, err)

	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, 0, len(engine.cw))
}

func TestMaxAttempts(t *testing.T) {
	t.Parallel()

	engine := NewEngine(WithMaxAttempts(2))
	assert.Equal(t, 2, engine.maxAttempts)

	engine.Start()
	defer engine.Stop()

	job := &jobExample{IsError: true}

	err := engine.Enqueue(context.TODO(), job)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(engine.cw))

	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, 2, job.AttemptsCount)
}

func TestMaxQueueSize(t *testing.T) {
	t.Parallel()

	engine := NewEngine(WithMaxQueueSize(1))
	assert.Equal(t, 1, engine.maxQueueSize)

	engine.Start()
	defer engine.Stop()

	sleepDuration := 2 * time.Second
	err := engine.Enqueue(context.TODO(), &jobExample{SleepDuration: &sleepDuration})
	assert.Nil(t, err)
	assert.Equal(t, 1, len(engine.cw))

	err = engine.Enqueue(context.TODO(), &jobExample{})
	assert.Equal(t, errors.New("The job's queue is full"), err)
}

func TestMaxWorkers(t *testing.T) {
	t.Parallel()

	engine := NewEngine(WithMaxWorkers(2))
	assert.Equal(t, 2, engine.maxWorkers)

	engine.Start()
	defer engine.Stop()

	err := engine.Enqueue(context.TODO(), &jobExample{})
	assert.Nil(t, err)
	assert.Equal(t, 1, len(engine.cw))

	err = engine.Enqueue(context.TODO(), &jobExample{})
	assert.Nil(t, err)
	assert.Equal(t, 2, len(engine.cw))

	err = engine.Enqueue(context.TODO(), &jobExample{})
	assert.Nil(t, err)
	assert.Equal(t, 2, len(engine.cw))
}

func TestInactivityTimeout(t *testing.T) {
	t.Parallel()

	engine := NewEngine(WithInactivityTimeout(10 * time.Millisecond))
	assert.Equal(t, 10*time.Millisecond, engine.inactivityTimeout)

	engine.Start()
	defer engine.Stop()

	err := engine.Enqueue(context.TODO(), &jobExample{})
	assert.Nil(t, err)
	assert.Equal(t, 1, len(engine.cw))

	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, 0, len(engine.cw))
}
