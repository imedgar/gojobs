package job

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Job struct {
	ID          uuid.UUID       `json:"id"`
	Type        string          `json:"type"`
	Payload     json.RawMessage `json:"payload"`
	Status      JobStatus       `json:"status"`
	Attempts    uint8           `json:"attempts"`
	MaxAttempts uint8           `json:"max_attempts"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

func NewJob(jobType string, payload json.RawMessage, maxAttempts uint8) *Job {
	now := time.Now().UTC()
	return &Job{
		ID:          uuid.New(),
		Type:        jobType,
		Payload:     payload,
		Status:      jobStatusValue["pending"],
		Attempts:    0,
		MaxAttempts: maxAttempts,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}
