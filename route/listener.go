package route

import "io"

/* Listener defines a type that can accept new connections
to receive requests */
type Listener interface {
	io.Closer
	Accept() (Conn, error)
}

/* Conn describes an object that can be used to read
request from */
type Conn interface {
	io.ReadWriteCloser
}
