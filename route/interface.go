package route

/* ReadRequest describes a function that can read a single
raw request from a Conn object */
type ReadRequest func(Conn) ([]byte, error)
