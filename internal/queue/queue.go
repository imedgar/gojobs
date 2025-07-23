package queue

import "gojobs/internal/job"

type Queue interface {
	Enqueue(*job.Job) error
	Dequeue() (*job.Job, error)
	Update(*job.Job) error
	GetByID(id string) (*job.Job, error)
}
