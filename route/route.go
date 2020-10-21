package route

import (
	"github.com/arstevens/go-request/handle"
)

/* UnpackAndRoute begins a pipeline for the listening, interpreting and routing
of requests from clients */
func UnpackAndRoute(listener Listener, done <-chan struct{}, handlers map[int]handle.RequestHandler,
	unpack handle.UnpackRequest, read ReadRequest) {
	identifyStream := make(chan RequestPair)
	pipelineDone := make(chan struct{})
	defer close(pipelineDone)

	go identifyAndRoute(identifyStream, handlers)
	go listenAndUnmarshal(listener, unpack, read, pipelineDone, identifyStream)
	<-done
}

// Route takes in a stream of handle.Requests and routes them to the appropriate handle.RequestHandler
func Route(requestStream <-chan RequestPair, done <-chan struct{}, handlers map[int]handle.RequestHandler) {
	go identifyAndRoute(requestStream, handlers)
	<-done
}
