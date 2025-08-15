package api

import (
	"bytes"
	"encoding/json"
	"gojobs/internal/job"
	"gojobs/internal/queue"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func decodeJobBody(t *testing.T, body io.Reader) *job.Job {
	t.Helper()
	var j job.Job
	err := json.NewDecoder(body).Decode(&j)
	if err != nil {
		t.Fatalf("failed to decode job: %v", err)
	}
	return &j
}

func TestHandler_CreateAndGetJob(t *testing.T) {
	q := queue.NewMemoryQueue(10)
	handler := NewHandler(q)

	server := httptest.NewServer(handler)
	defer server.Close()

	jobJSON := `{
		"type": "send_email",
		"payload": {"to": "user@example.com"},
		"max_attempts": 2
	}`

	resp, err := http.Post(server.URL+"/jobs", "application/json", bytes.NewBufferString(jobJSON))
	if err != nil {
		t.Fatalf("failed to create job: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		t.Fatalf("expected 202, got %d", resp.StatusCode)
	}

	created := decodeJobBody(t, resp.Body)

	if created.Type != "send_email" || created.MaxAttempts != 2 {
		t.Errorf("unexpected job fields: %+v", created)
	}

	getResp, err := http.Get(server.URL + "/jobs/" + created.ID.String())
	if err != nil {
		t.Fatalf("failed to get job: %v", err)
	}
	defer getResp.Body.Close()

	if getResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", getResp.StatusCode)
	}

	fetched := decodeJobBody(t, getResp.Body)
	if fetched.ID != created.ID {
		t.Errorf("job IDs don't match")
	}
}

func TestHandler_GetJob_InvalidID(t *testing.T) {
	q := queue.NewMemoryQueue(10)
	handler := NewHandler(q)

	req := httptest.NewRequest(http.MethodGet, "/jobs/not-a-uuid", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid UUID, got %d", rec.Code)
	}
}
