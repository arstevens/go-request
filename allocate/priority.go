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
func (pa *PriorityJobAllocator) AddJob(request interface{}, conn handle.Conn) error {
	var allocated bool
	startTime := time.Now()
	minIdx := 0
	minQueued := pa.handlers[minIdx].QueuedJobs()
	curIdx := nextIndex(minIdx, len(pa.handlers))
	for !allocated {
		queued := pa.handlers[curIdx].QueuedJobs()
		if queued < minQueued {
			minIdx = curIdx
			minQueued = queued
		}
		curIdx = nextIndex(curIdx, len(pa.handlers))
		if curIdx == 0 {
			if minQueued == pa.handlers[minIdx].JobCapacity() {
				if len(pa.handlers) == pa.handlerLimit {
					if time.Since(startTime) > AllocateTimeout {
						return AllocateLimitErr
					}
					time.Sleep(AllocateTimeout)
				} else {
					newHandler := pa.generator.NewHandler()
					newHandler.AddJob(request, conn)
					pa.handlers = append(pa.handlers, newHandler)
					allocated = true
				}
			} else {
				handler := pa.handlers[minIdx]
				handler.AddJob(request, conn)
				allocated = true
			}
		}
	}
	return nil
}

// QueuedJobs returns the number of queued jobs
func (pa *PriorityJobAllocator) QueuedJobs() int {
	queuedCount := 0
	for _, handler := range pa.handlers {
		queuedCount += handler.QueuedJobs()
	}
	return queuedCount
}

/* JobCapacity returns the max number of jobs that can
be queued at once */
func (pa *PriorityJobAllocator) JobCapacity() int {
	return pa.handlerLimit * pa.generator.HandlerCapacity()
}

// Close closes all the handlers
func (pa *PriorityJobAllocator) Close() error {
	var returnErr error
	for _, handler := range pa.handlers {
		err := handler.Close()
		if err != nil {
			returnErr = err
		}
	}
	return returnErr
}
