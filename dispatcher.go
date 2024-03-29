package go_run

type Dispatcher struct {
	// A pool of background_processing channels that are registered with the dispatcher
	maxWorkers int
	WorkerPool chan chan Job
	jobQueue   chan Job
	maxRetry   int
	logger     Logger
}

func NewDispatcher(config WorkerConfig, logger Logger) *Dispatcher {
	queue := make(chan Job, config.QueueSize)
	pool := make(chan chan Job, config.MaxWorkers)
	return &Dispatcher{
		WorkerPool: pool,
		maxWorkers: config.MaxWorkers,
		jobQueue:   queue,
		maxRetry:   config.MaxRetry,
		logger:     logger,
	}
}

func (d *Dispatcher) Run() {
	// starting n number of background_processing
	for i := 0; i < d.maxWorkers; i++ {
		worker := NewWorker(d.WorkerPool, d.maxRetry, d.logger)
		worker.Start()
	}

	go d.dispatch()
}

func (d *Dispatcher) Enqueue(job Job) {
	d.jobQueue <- job
}

func (d *Dispatcher) dispatch() {
	d.logger.Infof("Worker que dispatcher started...")
	for {

		select {
		case job := <-d.jobQueue:
			d.logger.Infof("a dispatcher request received")
			// a job request has been received
			go func(job Job) {
				// try to obtain a worker job channel that is available.
				// this will block until a worker is idle
				jobChannel := <-d.WorkerPool

				// dispatch the job to the worker job channel
				jobChannel <- job
			}(job)
		}
	}
}
