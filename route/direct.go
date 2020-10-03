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
func identifyAndRoute(requestStream <-chan handle.Request, handlers map[int]handle.RequestHandler) {
	for {
		request, ok := <-requestStream
		if !ok {
			log.Printf("Request Stream has been closed\n")
			return
		}

		initialType := request.GetType()
		handler, ok := handlers[initialType]
		if !ok {
			log.Println(UnknownRequestErr)
			continue
		}
		err := handler.AddJob(request)
		if err != nil {
			log.Println(err)
		}
	}
}
