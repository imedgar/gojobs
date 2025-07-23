package main

import (
	"encoding/json"
	"fmt"
	"gojobs/internal/job"
	"gojobs/internal/queue"
)

func main() {
	q := queue.NewMemoryQueue(10)

	payload := json.RawMessage(`{"task": "send_email"}`)
	j := job.NewJob("send_email", payload, 3)

	err := q.Enqueue(j)
	if err != nil {
		panic(err)
	}
	fmt.Println("Enqueued job ID:", j.ID)

	fetched, _ := q.Dequeue()
	fmt.Println("Dequeued job ID:", fetched.ID)
	fmt.Println("Status:", fetched.Status.String())

	fetched.Status = job.InProgress
	q.Update(fetched)

	byID, _ := q.GetByID(fetched.ID.String())
	fmt.Println("Updated status:", byID.Status.String())
}
