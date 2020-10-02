package handle

import "io"

/* RequestHandler describes an object that can handle a
stream of requests */
type RequestHandler interface {
	AddJob(interface{}) error
	JobCapacity() int
	QueuedJobs() int
	io.Closer
}

/* RequestHandlerGenerator generates new objects that
confrom to the RequestHandler interface */
type RequestHandlerGenerator interface {
	NewHandler() RequestHandler
}

/* Defines a function that can take a sequence of bytes
and attempt to unpack it into an object usable by a RequestHandler */
type UnpackRequest func([]byte) (interface{}, error)

/* Defines a function that can take an unmarshaled
request and return an identifying integer */
type GetRequestId func(interface{}) int
