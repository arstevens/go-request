package rtest

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/arstevens/go-request/allocate"
	"github.com/arstevens/go-request/handle"
	"github.com/arstevens/go-request/route"
)

func TestRequestLibrary(t *testing.T) {
	fmt.Printf("\nBASIC TEST\n========\n\n")
	requestCount := 30
	datapoints := make([]byte, requestCount)
	for i := 0; i < requestCount; i++ {
		datapoints[i] = byte(rand.Int())
	}

	listener := &TestListener{datapoints}
	done := make(chan struct{})

	gen := &TestGenerator{cap: 10}
	alloc := allocate.NewPriorityJobAllocator(10, gen)
	handlers := map[int]handle.RequestHandler{
		0: &TestHandler{capacity: 10},
		1: &TestHandler{capacity: 15},
		2: &TestHandler{capacity: 10},
		3: alloc,
	}

	go route.UnpackAndRoute(listener, done, handlers, UnpackTestRequest, ReadTestRequestFromConn)
	time.Sleep(time.Second * 5)
}
