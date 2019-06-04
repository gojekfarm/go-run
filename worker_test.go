package go_run

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestJob struct {
	resultChannel chan string
}

func (t TestJob) Execute() error {
	t.resultChannel <- "test"
	return nil
}

func TestWorker_Start(t *testing.T) {
	pool := make(chan chan Job, 1)
	jobChannel := make(chan Job, 1)
	w := Worker{
		WorkerPool: pool,
		JobChannel: jobChannel,
		quit:       make(chan bool, 1),
		logger:     TestLogger{},
	}
	w.Start()

	t.Run("registers worker to the worker pool", func(t *testing.T) {
		res := <-pool
		assert.Equal(t, res, w.JobChannel)
	})

	t.Run("Executes the job when it is pushed into the job channel", func(t *testing.T) {
		resultChannel := make(chan string, 1)
		testJob := TestJob{
			resultChannel: resultChannel,
		}
		jobChannel <- testJob
		res := <-resultChannel
		assert.Equal(t, res, "test")
	})
}

func TestWorker_Stop(t *testing.T) {
	t.Run("adds value to the quit channel", func(t *testing.T) {
		quitChannel := make(chan bool, 1)
		w := Worker{
			WorkerPool: make(chan chan Job, 1),
			JobChannel: make(chan Job, 1),
			quit:       quitChannel,
			logger:     TestLogger{},
		}
		w.Stop()
		res := <-quitChannel
		assert.Equal(t, res, true)
	})
}
