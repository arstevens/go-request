package route

import (
	"fmt"
	"log"

	"github.com/arstevens/go-request/handle"
)

// identifyAndRoute takes a request and sends it to the proper subcomponent
func identifyAndRoute(requestStream <-chan interface{}, getId handle.GetRequestId, handlers map[int]handle.RequestHandler) {
	for {
		request, ok := <-requestStream
		if !ok {
			return
		}

		initialType := getId(request)
		handler, ok := handlers[initialType]
		if !ok {
			log.Println(fmt.Errorf("Unknown code in identifyAndRoute()"))
			continue
		}
		err := handler.AddJob(request)
		if err != nil {
			log.Println(err)
		}
	}
}
