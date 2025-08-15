package worker

import (
	"gojobs/internal/job"
	"gojobs/internal/queue"
	"log"
	"math/rand"
	"time"
)

type WorkerPool struct {
	Queue       queue.Queue
	NumWorkers  int
	quitChannel chan struct{}
}

func NewWorkerPool(q queue.Queue, numWorkers int) *WorkerPool {
	return &WorkerPool{
		Queue:       q,
		NumWorkers:  numWorkers,
		quitChannel: make(chan struct{}),
	}
}

func (wp *WorkerPool) Start() {
	for i := 0; i < wp.NumWorkers; i++ {
		go wp.runWorker(i)
	}
	log.Printf("[WorkerPool] Started %d workers\n", wp.NumWorkers)
}

func (wp *WorkerPool) Stop() {
	close(wp.quitChannel)
	log.Println("[WorkerPool] Stopped")
}

func (wp *WorkerPool) runWorker(id int) {
	for {
		select {
		case <-wp.quitChannel:
			log.Printf("[Worker %d] Stopping\n", id)
			return
		default:
			jobItem, err := wp.Queue.Dequeue()
			if err != nil {
				time.Sleep(100 * time.Millisecond)
				continue
			}

			log.Printf("[Worker %d] Processing job %s (%s)\n", id, jobItem.ID, jobItem.Type)
			jobItem.Status = job.InProgress
			jobItem.Attempts++
			jobItem.UpdatedAt = time.Now()
			wp.Queue.Update(jobItem)

			success := wp.simulateJob(jobItem)
			if success {
				jobItem.Status = job.Done
				log.Printf("[Worker %d] Job %s done\n", id, jobItem.ID)
			} else {
				if jobItem.Attempts >= jobItem.MaxAttempts {
					jobItem.Status = job.Failed
					log.Printf("[Worker %d] Job %s failed permanently\n", id, jobItem.ID)
				} else {
					log.Printf("[Worker %d] Job %s failed - retrying\n", id, jobItem.ID)
					wp.Queue.Enqueue(jobItem)
					continue
				}
			}

			jobItem.UpdatedAt = time.Now()
			wp.Queue.Update(jobItem)
		}
	}
}

func (wp *WorkerPool) simulateJob(_ *job.Job) bool {
	time.Sleep(time.Duration(rand.Intn(500)+200) * time.Millisecond)
	return rand.Float32() < 0.8 // 80% success rate
}
