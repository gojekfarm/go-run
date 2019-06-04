package go_run

import "fmt"

type Worker struct {
	WorkerPool chan chan Job
	JobChannel chan Job
	MaxRetry   int
	quit       chan bool
}

func NewWorker(workerPool chan chan Job, maxRetry int) Worker {
	return Worker{
		WorkerPool: workerPool,
		JobChannel: make(chan Job),
		MaxRetry:   maxRetry,
		quit:       make(chan bool)}
}

// Start method starts the run loop for the worker, listening for a quit channel in
// case we need to stop it
func (w Worker) Start() {
	go func() {
		for {
			// register the current worker into the worker queue.
			w.WorkerPool <- w.JobChannel

			select {
			case job := <-w.JobChannel:
				fmt.Printf("Worker processing job")
				// we have received a work request.
				w.executeJob(job, 0)
			case <-w.quit:
				// we have received a signal to stop
				return
			}
		}
	}()
}

func (w Worker) executeJob(job Job, errorCount int) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Job panicked")
			w.retryIfErrorCount(job, errorCount)
		}
	}()

	if err := job.Execute(); err != nil {
		fmt.Printf("Error Executing Job: %s", err.Error())
		w.retryIfErrorCount(job, errorCount)
	} else {
		fmt.Printf("Worker finished processing job")
		return
	}
}

func (w Worker) retryIfErrorCount(job Job, errorCount int) {
	if errorCount < w.MaxRetry {
		fmt.Printf("Retrying Job")
		w.executeJob(job, errorCount+1)
	} else {
		fmt.Printf("Job exceeded retry count")
		return
	}
}

// Stop signals the worker to stop listening for work requests.
func (w Worker) Stop() {
	go func() {
		w.quit <- true
	}()
}
