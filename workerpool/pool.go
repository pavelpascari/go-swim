// Package workerpool is a Worker Pool abstraction built on top of channels with Synchronisation features.
package workerpool

// Job represents a struct to provide input data for the workers to process.
type Job struct {
	Args []any
}

// JobResult is a general struct that is passed as a result of a processed Job.
type JobResult struct {
	Result any
	Error  error
}

// WorkerFunc is used to inject custom logic of handling specific work.
//
// WorkerFunc is passed to the Run call.
type WorkerFunc func(Job) JobResult

// WorkerPool represents the API for a generic concurrent worker pool.
type WorkerPool interface {
	// Add is called to provide Jobs to the workers.
	// It is the responsibility of the caller to invoke `wp.Close()`
	// once all the needed work has been processed already.
	Add(Job)
	// Run starts the workers in separate goroutines and returns a <- chan JobResult.
	// It is the responsibility of the caller to process the results
	Run(w WorkerFunc) <-chan JobResult

	// Close is used to signal that all the needed work has been passed using the Add method and no additional Jobs are
	// expected.
	// Close will signal that the implementation can clean up the worker goroutines.
	// Note that Close is gracefully handling the signaling, and it might take a while
	// until all workers are through with their processing.
	Close() error

	// Stop is used to interrupt any ongoing processing.
	Stop()
}
