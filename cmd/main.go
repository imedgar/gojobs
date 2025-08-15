package main

import (
	"gojobs/internal/api"
	"gojobs/internal/queue"
	"gojobs/internal/worker"
	"log"
	"net/http"
)

func main() {
	q := queue.NewMemoryQueue(100)
	pool := worker.NewWorkerPool(q, 3)
	pool.Start()
	defer pool.Stop()

	handler := api.NewHandler(q)

	log.Println("[Main] Server listening on :8080")
	http.ListenAndServe(":8080", handler)
}
