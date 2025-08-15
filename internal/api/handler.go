package api

import (
	"encoding/json"
	"gojobs/internal/job"
	"gojobs/internal/queue"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type Handler struct {
	Queue queue.Queue
}

func NewHandler(q queue.Queue) *Handler {
	return &Handler{Queue: q}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodPost && r.URL.Path == "/jobs":
		h.handleCreateJob(w, r)
	case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/jobs/"):
		h.handleGetJob(w, r)
	default:
		http.NotFound(w, r)
	}
}

type createJobRequest struct {
	Type        string          `json:"type"`
	Payload     json.RawMessage `json:"payload"`
	MaxAttempts uint8           `json:"max_attempts"`
}

func (h *Handler) handleCreateJob(w http.ResponseWriter, r *http.Request) {
	var req createJobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	j := job.NewJob(req.Type, req.Payload, req.MaxAttempts)
	if err := h.Queue.Enqueue(j); err != nil {
		http.Error(w, "failed to enqueue job", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(j)
}

func (h *Handler) handleGetJob(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 || parts[1] != "jobs" {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	idStr := parts[2]
	if _, err := uuid.Parse(idStr); err != nil {
		http.Error(w, "invalid UUID", http.StatusBadRequest)
		return
	}

	j, err := h.Queue.GetByID(idStr)
	if err != nil {
		http.Error(w, "job not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(j)
}
