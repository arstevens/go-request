package route

import (
	"github.com/arstevens/go-request/handle"
)

/* BeginRouting begins a pipeline for the listening, interpreting and routing
of requests from clients */
func BeginRouting(listener Listener, done <-chan struct{}, handlers map[int]handle.RequestHandler,
	getId handle.GetRequestId, unpack handle.UnpackRequest, read ReadRequest) {
	identifyStream := make(chan interface{})
	pipelineDone := make(chan struct{})
	defer close(pipelineDone)

	go identifyAndRoute(identifyStream, getId, handlers)
	go listenAndUnmarshal(listener, unpack, read, pipelineDone, identifyStream)
	<-done
}
