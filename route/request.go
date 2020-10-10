package route

import (
	"github.com/arstevens/go-request/handle"
)

/* Defines a function that can take a sequence of bytes
and attempt to unpack it into an object usable by a RequestHandler */
type UnpackRequest func([]byte, Conn) (handle.Request, error)
