package route

import (
	"errors"
	"log"

	"github.com/arstevens/go-request/handle"
)

/* UnknownRequestErr is an error indicating receiving a request
that could not be mapped to a handler */
var UnknownRequestErr = errors.New("Unknown request type code")

// identifyAndRoute takes a request and sends it to the proper subcomponent
func identifyAndRoute(requestStream <-chan handle.RequestPair, handlers map[int]handle.RequestHandler) {
	for {
		requestPair, ok := <-requestStream
		if !ok {
			log.Printf("Request Stream has been closed\n")
			return
		}

		request := requestPair.Request.(handle.Request)
		initialType := request.GetType()
		handler, ok := handlers[initialType]
		if !ok {
			log.Println(UnknownRequestErr)
			continue
		}
		err := handler.AddJob(requestPair.Request, requestPair.Conn)
		if err != nil {
			log.Println(err)
		}
	}
}
