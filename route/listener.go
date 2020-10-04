package route

import (
	"io"
	"net"
)

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

// NetListener wraps a net.Listener so it implements the Listener interface
type NetListener struct {
	listener net.Listener
}

// Accept returns a net.Conn wrapped as a NetConn and an error
func (nl *NetListener) Accept() (Conn, error) {
	nConn, err := nl.listener.Accept()
	conn := NetConn{conn: nConn}
	return &conn, err
}

// Close closes the underlying net.Listener
func (nl *NetListener) Close() error {
	return nl.listener.Close()
}

// NetConn wraps a net.Conn so it implements the Conn interface
type NetConn struct {
	conn net.Conn
}

// Read reads data from the net.Conn into a slice of bytes
func (nc *NetConn) Read(b []byte) (int, error) {
	return nc.conn.Read(b)
}

// Write writes data from a slice of bytes to the underlying net.Conn
func (nc *NetConn) Write(b []byte) (int, error) {
	return nc.conn.Write(b)
}

// Close closes the underlying net.Conn
func (nc *NetConn) Close() error {
	return nc.conn.Close()
}
