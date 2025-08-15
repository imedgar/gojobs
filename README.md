# GoJobs

A minimal, event-driven job queue system in Go

---

## Features

- In-memory job queue with channel-based FIFO
- Worker pool with concurrent job execution
- Retry logic with max attempt control
- HTTP API to enqueue and track jobs

---

## Project Structure

```
gojobs/
├── cmd/              # Entry point
├── internal/
│   ├── api/          # HTTP handlers (net/http)
│   ├── job/          # Job model & statuses
│   ├── queue/        # In-memory queue implementation
│   └── worker/       # Worker pool logic
```

---

## Getting Started

### 1. Clone and run

```bash
git clone https://github.com/imedgar/gojobs.git
cd gojobs
go run ./...
```

This will:

- Start the HTTP API at `http://localhost:8080`
- Start the worker pool with 3 workers

---

### 2. Enqueue a job

```bash
curl -X POST http://localhost:8080/jobs \
  -H "Content-Type: application/json" \
  -d '{
    "type": "send_email",
    "payload": {"to": "user@example.com"},
    "max_attempts": 3
  }'
```

### 3. Get job status

```bash
curl http://localhost:8080/jobs/{job_id}
```

---

## 🧪 Running Tests

```bash
go test ./...
```

Includes unit and integration tests for:

- Memory queue
- Worker pool
- HTTP API

---

## Design Decisions

- All concurrency is handled using `goroutines` and `channels`
- `map[string]*Job` + `sync.Mutex` used for state tracking
- Clean separation of concerns via `internal/` packages

---

## Future Ideas

- [ ] Add persistent queue with file or DB
- [ ] Support job scheduling (delayed jobs)
- [ ] Add priorities or tags
- [ ] Web dashboard or CLI interface
- [ ] Redis-backed distributed queue (plugged into same interface)

---

## Author

[@imedgar](https://github.com/imedgar)

---

## License

MIT
