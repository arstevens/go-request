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

func TestPipelining(t *testing.T) {
	fmt.Printf("\nPIPELINE TEST\n=============\n\n")
	requestCount := 30
	datapoints := make([]byte, requestCount)
	for i := 0; i < requestCount; i++ {
		datapoints[i] = byte(rand.Int())
	}

	listener := &TestListener{datapoints}
	newConverge := make(chan (<-chan handle.RequestPair), 1)
	gen := &TestPipeGenerator{cap: 5, newOut: newConverge}
	alloc := allocate.NewCyclicJobAllocator(10, gen)

	handlerCount := 3
	chans := make([]chan handle.RequestPair, handlerCount)
	handlers := map[int]handle.RequestHandler{3: alloc}
	for i := 0; i < handlerCount; i++ {
		chans[i] = make(chan handle.RequestPair)
		handlers[i] = &TestPipeHandler{capacity: 10, out: chans[i]}
	}

	endpointHandlers := map[int]handle.RequestHandler{
		0: &TestHandler{capacity: 10},
		1: &TestHandler{capacity: 12},
		2: &TestHandler{capacity: 8},
		3: &TestHandler{capacity: 20},
	}

	pipe := make(chan handle.RequestPair)
	go route.ConvergeChannels([]<-chan handle.RequestPair{chans[0], chans[1],
		chans[2]}, newConverge, pipe)

	done := make(chan struct{})
	defer close(done)
	go route.Route(pipe, done, endpointHandlers)
	go route.UnpackAndRoute(listener, done, handlers, UnpackTestRequest, ReadTestRequestFromConn)

	time.Sleep(time.Second * 2)
	for _, c := range chans {
		close(c)
	}
	time.Sleep(time.Second * 3)
}
