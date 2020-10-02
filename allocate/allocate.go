package allocate

import (
	"errors"

	"github.com/arstevens/go-request/handle"
)

var AllocateLimitErr = errors.New("Job Allocation Limit Reached")

/* CyclicJobAllocator is an implementation of RequestHandler
the allocates jobs by cycling through handlers and allocating
new handlers when all are backed up */
type CyclicJobAllocator struct {
	handlers     []handle.RequestHandler
	generator    handle.RequestHandlerGenerator
	index        int
	handlerLimit int
}

/* NewCyclicJobAllocator creates a new CyclicJobAllocator that uses the
provided generator to create new transaction handlers when needed */
func NewCyclicJobAllocator(handlerLimit int, generator handle.RequestHandlerGenerator) *CyclicJobAllocator {
	alloc := &CyclicJobAllocator{
		handlers:     make([]handle.RequestHandler, 0),
		generator:    generator,
		index:        0,
		handlerLimit: handlerLimit,
	}
	alloc.handlers = append(alloc.handlers, generator.NewHandler())
	return alloc
}

// AllocateJob allocates a request to a TransactionHandler
func (ca *CyclicJobAllocator) AddJob(request interface{}) error {
	var allocated bool
	startIdx := ca.index
	for !allocated {
		handler := ca.handlers[ca.index]
		if handler.QueuedJobs() != handler.JobCapacity() {
			handler.AddJob(request)
			allocated = true
		} else {
			ca.index = (ca.index + 1) % len(ca.handlers)
			if ca.index == startIdx {
				if len(ca.handlers) == ca.handlerLimit {
					return AllocateLimitErr
				}
				newHandler := ca.generator.NewHandler()
				newHandler.AddJob(request)
				ca.handlers = append(ca.handlers, newHandler)
				ca.index = 0
				allocated = true
			}
		}
	}
	return nil
}

// QueuedJobs returns the number of queued jobs
func (ca *CyclicJobAllocator) QueuedJobs() int {
	queuedCount := 0
	for _, handler := range ca.handlers {
		queuedCount += handler.QueuedJobs()
	}
	return queuedCount
}

// Close closes the cyclic allocator
func (ca *CyclicJobAllocator) Close() {
	for _, handler := range ca.handlers {
		handler.Close()
	}
}
