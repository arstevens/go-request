package rtest

import (
	"errors"
	"fmt"

	"github.com/arstevens/go-request/handle"
)

type TestPipeHandler struct {
	capacity int
	queued   int
	out      chan handle.RequestPair
}

func (h *TestPipeHandler) AddJob(i interface{}, c handle.Conn) error {
	h.out <- handle.RequestPair{i.(handle.Request), c}
	return nil
}

func (h *TestPipeHandler) JobCapacity() int {
	return h.capacity
}

func (h *TestPipeHandler) QueuedJobs() int {
	return 0 //h.capacity
}

func (h *TestPipeHandler) Close() error {
	return nil
}

type TestPipeGenerator struct {
	cap    int
	newOut chan<- <-chan handle.RequestPair
}

func (g *TestPipeGenerator) NewHandler() handle.RequestHandler {
	nChan := make(chan handle.RequestPair)
	g.newOut <- nChan
	return &TestPipeHandler{capacity: g.cap, out: nChan}
}

func (g *TestPipeGenerator) HandlerCapacity() int {
	return g.cap
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

func (t *TestListener) Accept() (handle.Conn, error) {
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

type TestRequest struct {
	code    int32
	request []byte
}

func (t *TestRequest) GetType() int32 {
	return t.code
}

func (t *TestRequest) GetRequest() []byte {
	return t.request
}

func ReadTestRequestFromConn(c handle.Conn) ([]byte, error) {
	b := make([]byte, 1)
	c.Read(b)
	return []byte{b[0] % 4}, nil
}

func UnpackTestWrapperRequest(b []byte) (handle.Request, error) {
	x := TestRequest{
		code:    int32(b[0]),
		request: []byte{},
	}
	return &x, nil
}

func UnpackTestRequest(b []byte) (interface{}, error) {
	return nil, nil
}

type TestHandler struct {
	capacity int
	queued   int
}

func (h *TestHandler) AddJob(interface{}, handle.Conn) error {
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
