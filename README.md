# Gawe

[![Build Status](https://github.com/nadjieb/gawe/workflows/Build/badge.svg)](https://github.com/nadjieb/gawe/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/nadjieb/gawe)](https://goreportcard.com/report/github.com/nadjieb/gawe)
[![Maintainability](https://api.codeclimate.com/v1/badges/c3c92fbf37c8e26281b4/maintainability)](https://codeclimate.com/github/nadjieb/gawe/maintainability)
[![Codecov](https://codecov.io/gh/nadjieb/gawe/branch/master/graph/badge.svg)](https://codecov.io/gh/nadjieb/gawe)

## Description
Gawe is a Go library for processing background jobs using Go channels as FIFO queue to control job execution and worker instantiation.

## Installation
```sh
go get -u github.com/nadjieb/gawe
```

## Usage
```go
package main

import (
	"context"
	"time"

	"github.com/nadjieb/gawe"
)

// RecordHistoryJob is a struct that comply to gawe.Job interface
type RecordHistoryJob struct {
	ID   string
	Data string
}

var _ gawe.Job = (*RecordHistoryJob)(nil)

// JobID returns ID of the job (Usually used for logging)
func (j *RecordHistoryJob) JobID() string {
	return j.ID
}

// JobType returns type of the job (Usually used for logging)
func (j *RecordHistoryJob) JobType() string {
	return "record-history"
}

// Tags returns tags of the job (Usually used for logging)
func (j *RecordHistoryJob) Tags() []string {
	return []string{"record", "history"}
}

// Exec execute the job
func (j *RecordHistoryJob) Exec(ctx context.Context) error {
	var err error

	// record history

	return err
}

func main() {
	engine := gawe.NewEngine(
		gawe.WithMaxAttempts(3),                   // max attempts of job executions if failed
		gawe.WithMaxQueueSize(100),                // max queue size for jobs
		gawe.WithMaxWorkers(4),                    // max workers run in the background
		gawe.WithInactivityTimeout(5*time.Second), // a worker will stop running since last defined inactivity timeout after last job execution
	)

	engine.Start()

	job := &RecordHistoryJob{ID: "123abc", Data: "record"}

	err := engine.Enqueue(context.Background(), job)
	if err != nil {
		// handle error
	}

	engine.Stop()
}
```

### Plugins
To create a plugin for the engine, create a struct that fulfill the [Plugin](plugin.go) interface then add it to gawe engine as an [Option](option.go).

```go
// Logger is a struct that comply to gawe.Plugin interface
type Logger struct{}

var _ gawe.Plugin = (*Logger)(nil)

// OnJobStart is called just before the job execution
func (l *Logger) OnJobStart(ctx context.Context, job gawe.IdentifiableJob) context.Context {
	// return the (new) context to pass it to the next plugin/job
	return ctx
}

// OnJobEnd is called once the job has successfully executed
func (l *Logger) OnJobEnd(ctx context.Context, job gawe.IdentifiableJob) {
	// do stuffs
}

// OnJobError is called if the job execution failed
func (l *Logger) OnJobError(ctx context.Context, job gawe.IdentifiableJob, err error) context.Context {
	// return the (new) context to pass it to the next plugin/job
	return ctx
}

...

logger := &Logger{}
engine := gawe.NewEngine(gawe.WithPlugins(logger))
```

## License
Released under the [Apache License 2.0](LICENSE)
