package worker

import (
	"fmt"
	"time"
)

type Job struct {
	UserID  string
	EventID string
}

type Pool struct {
	JobQueue chan Job
}

func NewPool(numWorkers int, queueSize int) *Pool {
	pool := &Pool{
		JobQueue: make(chan Job, queueSize),
	}
	
	for i := 1; i <= numWorkers; i++ {
		go pool.worker(i)
	}

	return pool
}

func (p *Pool) worker(id int) {
	for job := range p.JobQueue {
		time.Sleep(1 * time.Second)
		fmt.Printf("[Worker %d] Confirmation sent to User %s for Event %s\n", id, job.UserID, job.EventID)
	}
}

func (p *Pool) Enqueue(job Job) {
	p.JobQueue <- job
}
