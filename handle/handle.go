package handle

/* RequestHandler describes an object that can handle a
stream of requests */
type RequestHandler interface {
	AddJob(interface{}) error
	JobCapacity() int
	QueuedJobs() int
}

/* RequestHandlerGenerator generates new objects that
confrom to the RequestHandler interface */
type RequestHandlerGenerator interface {
	NewHandler() RequestHandler
}

/* Defines a Request object with only one method to
retrieve an integer identifying the type of the request */
type Request interface {
	GetType() int
}

/* Defines a function that can take a sequence of bytes
and attempt to unpack it into an object usable by a RequestHandler */
type UnpackRequest func([]byte) (Request, error)
