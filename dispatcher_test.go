package go_run

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDispatcher(t *testing.T) {
	resultChannel := make(chan string, 1)
	testJob := TestJob{
		resultChannel: resultChannel,
	}

	d := NewDispatcher(
		WorkerConfig{
			MaxWorkers: 2,
			QueueSize:  5,
			MaxRetry:   1,
		},
		TestLogger{},
	)
	d.Run()

	t.Run("performs the job enqueued", func(t *testing.T) {
		d.Enqueue(testJob)
		res := <-resultChannel
		assert.Equal(t, res, "test")
	})
}
