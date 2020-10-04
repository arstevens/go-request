package allocate

import (
	"time"

	"github.com/arstevens/go-request/handle"
)

/* PriorityJobAllocator is an implementation of RequestHandler
the allocates jobs by choosing the handler with the least number
of queued jobs and allocating new handlers when all queued are filled */
type PriorityJobAllocator struct {
	handlers     []handle.RequestHandler
	generator    handle.RequestHandlerGenerator
	handlerLimit int
}

/* NewPriorityJobAllocator creates a new PriorityJobAllocator that uses the
provided generator to create new request handlers when needed */
func NewPriorityJobAllocator(handlerLimit int, generator handle.RequestHandlerGenerator) *PriorityJobAllocator {
	alloc := &PriorityJobAllocator{
		handlers:     make([]handle.RequestHandler, 0),
		generator:    generator,
		handlerLimit: handlerLimit,
	}
	alloc.handlers = append(alloc.handlers, alloc.generator.NewHandler())
	return alloc
}

// AddJob allocates a job to a RequestHandler
func (ca *PriorityJobAllocator) AddJob(request interface{}) error {
	var allocated bool
	startTime := time.Now()
	minIdx := 0
	minQueued := ca.handlers[minIdx].QueuedJobs()
	curIdx := nextIndex(minIdx, len(ca.handlers))
	for !allocated {
		queued := ca.handlers[curIdx].QueuedJobs()
		if queued < minQueued {
			minIdx = curIdx
			minQueued = queued
		}
		curIdx = nextIndex(curIdx, len(ca.handlers))
		if curIdx == 0 {
			if minQueued == ca.handlers[minIdx].JobCapacity() {
				if len(ca.handlers) == ca.handlerLimit {
					if time.Since(startTime) > AllocateTimeout {
						return AllocateLimitErr
					}
					time.Sleep(AllocateTimeout)
				} else {
					newHandler := ca.generator.NewHandler()
					newHandler.AddJob(request)
					ca.handlers = append(ca.handlers, newHandler)
					allocated = true
				}
			} else {
				handler := ca.handlers[minIdx]
				handler.AddJob(request)
				allocated = true
			}
		}
	}
	return nil
}

// QueuedJobs returns the number of queued jobs
func (ca *PriorityJobAllocator) QueuedJobs() int {
	queuedCount := 0
	for _, handler := range ca.handlers {
		queuedCount += handler.QueuedJobs()
	}
	return queuedCount
}

/* JobCapacity returns the max number of jobs that can
be queued at once */
func (ca *PriorityJobAllocator) JobCapacity() int {
	return ca.handlerLimit * ca.generator.HandlerCapacity()
}
