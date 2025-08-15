package queue

import (
	"errors"
	"gojobs/internal/job"
	"sync"
)

type MemoryQueue struct {
	JobsChan chan *job.Job
	JobsMap  map[string]*job.Job
	Mutex    sync.RWMutex
}

func NewMemoryQueue(bufferSize int) *MemoryQueue {
	return &MemoryQueue{
		JobsChan: make(chan *job.Job, bufferSize),
		JobsMap:  make(map[string]*job.Job),
	}
}

func (q *MemoryQueue) Enqueue(j *job.Job) error {
	q.Mutex.Lock()
	defer q.Mutex.Unlock()

	q.JobsMap[j.ID.String()] = j
	q.JobsChan <- j
	return nil
}

func (q *MemoryQueue) Dequeue() (*job.Job, error) {
	j, ok := <-q.JobsChan
	if !ok {
		return nil, errors.New("queue closed")
	}
	return j, nil
}

func (q *MemoryQueue) Update(j *job.Job) error {
	q.Mutex.Lock()
	defer q.Mutex.Unlock()

	if _, exists := q.JobsMap[j.ID.String()]; !exists {
		return errors.New("job not found")
	}

	q.JobsMap[j.ID.String()] = j
	return nil
}

func (q *MemoryQueue) GetByID(id string) (*job.Job, error) {
	q.Mutex.RLock()
	defer q.Mutex.RUnlock()

	if j, exists := q.JobsMap[id]; exists {
		return j, nil
	}

	return nil, errors.New("job not found")
}
