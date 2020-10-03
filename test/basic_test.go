package rtest

import (
	"errors"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/arstevens/go-request/allocate"
	"github.com/arstevens/go-request/handle"
	"github.com/arstevens/go-request/route"
)

func TestRequestLibrary(t *testing.T) {
	requestCount := 100
	datapoints := make([]byte, requestCount)
	for i := 0; i < requestCount; i++ {
		datapoints[i] = byte(rand.Int())
	}

	listener := &TestListener{datapoints}
	done := make(chan struct{})

	gen := &TestGenerator{cap: 3}
	alloc := allocate.NewCyclicJobAllocator(10, gen)
	handlers := map[int]handle.RequestHandler{
		0: &TestHandler{capacity: 10},
		1: &TestHandler{capacity: 15},
		2: &TestHandler{capacity: 10},
		3: alloc,
	}

	go route.UnpackAndRoute(listener, done, handlers, UnpackTestRequest, ReadTestRequestFromConn)
	time.Sleep(time.Second * 5)
}

type TestGenerator struct {
	cap int
}

func (g *TestGenerator) NewHandler() handle.RequestHandler {
	return &TestHandler{capacity: g.cap}
}

func (g *TestGenerator) HandlerCapacity() int {
	return g.cap
}

type TestListener struct {
	datapoints []byte
}

func (t *TestListener) Accept() (route.Conn, error) {
	if len(t.datapoints) == 0 {
		return nil, errors.New("out of stuff")
	}
	point := t.datapoints[0]
	t.datapoints = t.datapoints[1:]
	return &TestConn{point}, nil
}

func (t *TestListener) Close() error {
	return nil
}

type TestConn struct {
	data byte
}

func (t *TestConn) Read(b []byte) (int, error) {
	b[0] = t.data
	return 0, nil
}

func (t *TestConn) Write(b []byte) (int, error) {
	return 0, nil
}

func (t *TestConn) Close() error {
	return nil
}

type TestRequest int

func (t *TestRequest) GetType() int {
	return int(*t)
}

func ReadTestRequestFromConn(c route.Conn) ([]byte, error) {
	b := make([]byte, 1)
	c.Read(b)
	return []byte{b[0] % 4}, nil
}

func UnpackTestRequest(b []byte) (handle.Request, error) {
	x := TestRequest(b[0])
	return &x, nil
}

type TestHandler struct {
	capacity int
	queued   int
}

func (h *TestHandler) AddJob(interface{}) error {
	if h.queued != h.capacity {
		h.queued++
	}
	fmt.Printf("Received new job! %d/%d spots used\n", h.queued, h.capacity)
	return nil
}

func (h *TestHandler) JobCapacity() int {
	return h.capacity
}

func (h *TestHandler) QueuedJobs() int {
	return h.queued
}

func (h *TestHandler) Close() error {
	return nil
}
