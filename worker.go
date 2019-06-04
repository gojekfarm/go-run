package go_run

type Worker struct {
	WorkerPool chan chan Job
	JobChannel chan Job
	MaxRetry   int
	quit       chan bool
	logger     Logger
}

func NewWorker(workerPool chan chan Job, maxRetry int, logger Logger) Worker {
	return Worker{
		WorkerPool: workerPool,
		JobChannel: make(chan Job),
		MaxRetry:   maxRetry,
		quit:       make(chan bool),
		logger:     logger,
	}
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
				w.logger.Infof("Worker processing job")
				// we have received a work request.
				w.executeJob(job, 0)
			case <-w.quit:
				// we have received a signal to stop
				w.logger.Infof("Worker shutting down")
				return
			}
		}
	}()
}

func (w Worker) executeJob(job Job, errorCount int) {
	defer func() {
		if r := recover(); r != nil {
			w.logger.Errorf("Job panicked")
			w.retryIfErrorCount(job, errorCount)
		}
	}()

	if err := job.Execute(); err != nil {
		w.logger.Errorf("Error Executing Job: %s", err.Error())
		w.retryIfErrorCount(job, errorCount)
	} else {
		w.logger.Infof("Worker finished processing job")
		return
	}
}

func (w Worker) retryIfErrorCount(job Job, errorCount int) {
	if errorCount < w.MaxRetry {
		w.logger.Errorf("Retrying Job")
		w.executeJob(job, errorCount+1)
	} else {
		w.logger.Errorf("Job exceeded retry count")
		return
	}
}

// Stop signals the worker to stop listening for work requests.
func (w Worker) Stop() {
	go func() {
		w.quit <- true
	}()
}
