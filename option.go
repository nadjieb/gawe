package gawe

import "time"

// Option represents the workers controller options
type Option func(*Engine)

// WithMaxAttempts sets the max attempts of a job execution
func WithMaxAttempts(maxAttempts int) Option {
	return func(e *Engine) {
		e.maxAttempts = maxAttempts
	}
}

// WithMaxQueueSize sets the max queue of the jobs to be execute by the workers
func WithMaxQueueSize(maxQueueSize int) Option {
	return func(e *Engine) {
		e.maxQueueSize = maxQueueSize
	}
}

// WithMaxWorkers sets the max workers can be run in the background
func WithMaxWorkers(maxWorkers int) Option {
	return func(e *Engine) {
		e.maxWorkers = maxWorkers
	}
}

// WithInactivityTimeout sets the lifetime timeout of a worker since the last job execution
func WithInactivityTimeout(inactivityTimeout time.Duration) Option {
	return func(e *Engine) {
		e.inactivityTimeout = inactivityTimeout
	}
}

// WithPlugins sets the plugins
func WithPlugins(plugins ...Plugin) Option {
	return func(e *Engine) {
		for _, p := range plugins {
			e.plugins = append(e.plugins, p)
		}
	}
}
