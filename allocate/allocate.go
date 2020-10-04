package allocate

import (
	"errors"
	"time"

	"github.com/arstevens/go-request/handle"
)

var AllocateLimitErr = errors.New("Job Allocation Limit Reached")
var AllocateTimeout = time.Second

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
	alloc.handlers = append(alloc.handlers, alloc.generator.NewHandler())
	return alloc
}

// AllocateJob allocates a request to a TransactionHandler
func (ca *CyclicJobAllocator) AddJob(request interface{}) error {
	var allocated bool
	startIdx := ca.index
	startTime := time.Now()
	for !allocated {
		handler := ca.handlers[ca.index]
		if handler.QueuedJobs() != handler.JobCapacity() {
			handler.AddJob(request)
			ca.index = nextIndex(ca.index, len(ca.handlers))
			allocated = true
		} else {
			ca.index = nextIndex(ca.index, len(ca.handlers))
			if ca.index == startIdx {
				if len(ca.handlers) == ca.handlerLimit {
					if time.Since(startTime) > AllocateTimeout {
						return AllocateLimitErr
					}
					ca.index = nextIndex(ca.index, len(ca.handlers))
					time.Sleep(AllocateTimeout)
				} else {
					newHandler := ca.generator.NewHandler()
					newHandler.AddJob(request)
					ca.handlers = append(ca.handlers, newHandler)
					ca.index = 0
					allocated = true
				}
			}
		}
	}
	return nil
}

func nextIndex(curIdx int, maxIdx int) int {
	return (curIdx + 1) % maxIdx
}

// QueuedJobs returns the number of queued jobs
func (ca *CyclicJobAllocator) QueuedJobs() int {
	queuedCount := 0
	for _, handler := range ca.handlers {
		queuedCount += handler.QueuedJobs()
	}
	return queuedCount
}

/* JobCapacity returns the max number of jobs that can
be queued at once */
func (ca *CyclicJobAllocator) JobCapacity() int {
	return ca.handlerLimit * ca.generator.HandlerCapacity()
}
