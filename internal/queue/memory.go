package queue

import (
	"errors"
	"gojobs/internal/job"
	"sync"
)

type MemoryQueue struct {
	jobsChan chan *job.Job
	jobsMap  map[string]*job.Job
	mutex    sync.RWMutex
}

func NewMemoryQueue(bufferSize int) *MemoryQueue {
	return &MemoryQueue{
		jobsChan: make(chan *job.Job, bufferSize),
		jobsMap:  make(map[string]*job.Job),
	}
}

func (q *MemoryQueue) Enqueue(j *job.Job) error {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	q.jobsMap[j.ID.String()] = j
	q.jobsChan <- j
	return nil
}

func (q *MemoryQueue) Dequeue() (*job.Job, error) {
	j, ok := <-q.jobsChan
	if !ok {
		return nil, errors.New("queue closed")
	}
	return j, nil
}

func (q *MemoryQueue) Update(j *job.Job) error {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if _, exists := q.jobsMap[j.ID.String()]; !exists {
		return errors.New("job not found")
	}

	q.jobsMap[j.ID.String()] = j
	return nil
}

func (q *MemoryQueue) GetByID(id string) (*job.Job, error) {
	q.mutex.RLock()
	defer q.mutex.RUnlock()

	if j, exists := q.jobsMap[id]; exists {
		return j, nil
	}

	return nil, errors.New("job not found")
}
