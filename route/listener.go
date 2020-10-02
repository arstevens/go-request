package route

import "io"

type Listener interface {
	io.Closer
	Accept() (Conn, error)
}

type Conn interface {
	io.ReadWriteCloser
}
