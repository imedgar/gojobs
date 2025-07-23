package queue_test

import (
	"encoding/json"
	"gojobs/internal/job"
	"gojobs/internal/queue"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestMemoryQueue_BasicFlow(t *testing.T) {
	q := queue.NewMemoryQueue(10)

	payload := json.RawMessage(`{"key": "value"}`)
	newJob := job.NewJob("test_type", payload, 3)

	err := q.Enqueue(newJob)
	assert.NoError(t, err)

	dequeuedJob, err := q.Dequeue()
	assert.NoError(t, err)
	assert.Equal(t, newJob.ID, dequeuedJob.ID)
	assert.Equal(t, job.Pending, dequeuedJob.Status)

	dequeuedJob.Status = job.InProgress
	err = q.Update(dequeuedJob)
	assert.NoError(t, err)

	found, err := q.GetByID(dequeuedJob.ID.String())
	assert.NoError(t, err)
	assert.Equal(t, job.InProgress, found.Status)
}

func TestMemoryQueue_GetByID_NotFound(t *testing.T) {
	q := queue.NewMemoryQueue(10)

	_, err := q.GetByID(uuid.New().String())
	assert.Error(t, err)
}

func TestMemoryQueue_Update_NonExistentJob(t *testing.T) {
	q := queue.NewMemoryQueue(10)

	fakeJob := job.NewJob("fake", json.RawMessage(`{}`), 1)
	err := q.Update(fakeJob)
	assert.Error(t, err)
}
