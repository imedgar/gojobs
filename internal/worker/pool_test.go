package worker_test

import (
	"encoding/json"
	"gojobs/internal/job"
	"gojobs/internal/queue"
	"gojobs/internal/worker"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWorkerPool_ProcessesJobs(t *testing.T) {
	q := queue.NewMemoryQueue(10)
	pool := worker.NewWorkerPool(q, 2)
	pool.Start()
	defer pool.Stop()

	var jobIDs []string
	for i := 0; i < 3; i++ {
		payload := json.RawMessage(`{"type": "test"}`)
		j := job.NewJob("test_type", payload, 2)
		jobIDs = append(jobIDs, j.ID.String())
		err := q.Enqueue(j)
		assert.NoError(t, err)
	}

	waitFor(t, 5*time.Second, 100*time.Millisecond, func() bool {
		for _, id := range jobIDs {
			j, err := q.GetByID(id)
			if err != nil {
				return false
			}
			if j.Status != job.Done && j.Status != job.Failed {
				return false
			}
		}
		return true
	})

	for _, id := range jobIDs {
		j, err := q.GetByID(id)
		assert.NoError(t, err)
		assert.Contains(t, []job.JobStatus{job.Done, job.Failed}, j.Status)
	}
}

func waitFor(t *testing.T, timeout, interval time.Duration, condition func() bool) {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if condition() {
			return
		}
		time.Sleep(interval)
	}
	t.Fatal("Timeout")
}
