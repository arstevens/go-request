package route

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"

	"github.com/arstevens/go-request/handle"
)

/* RequestPair is a datatype composed of the received request
and a connection to the party that sent the request */
type RequestPair struct {
	Request handle.Request
	Conn    handle.Conn
}

/* listenAndUnmarshal accepts any connections from listener and attempts
to read and deserialize a request */
func listenAndUnmarshal(listener Listener, unpacker handle.UnpackRequest, reader ReadRequest,
	done <-chan struct{}, outStream chan<- RequestPair) {
	defer close(outStream)
	defer listener.Close()

	requestChan := make(chan RequestPair)
	go receiveRequests(listener, unpacker, reader, requestChan)
	for {
		select {
		case request, ok := <-requestChan:
			if !ok {
				return
			}
			outStream <- request
		case <-done:
			return
		}
	}
}

/* receiveRequests accepts all connections on the listener and
attempts to deserialize them into handle.Request objects. It then passes
these objects through the returnStream channel */
func receiveRequests(listener Listener, unpacker handle.UnpackRequest, reader ReadRequest, returnStream chan<- RequestPair) {
	defer close(returnStream)
	for {
		conn, err := listener.Accept()
		/* An error will be returned if the listener was closed.
		this allows gracefully stopping */
		if err != nil {
			log.Println(err)
			return
		}

		request, err := readAndUnpackRequest(conn, reader, unpacker)
		if err != nil {
			conn.Close()
			log.Println(err)
			continue
		}
		returnStream <- RequestPair{request, conn}
	}
}

// Consolidates the reading of a request
func readAndUnpackRequest(conn handle.Conn, reader ReadRequest, unpacker handle.UnpackRequest) (handle.Request, error) {
	rawRequest, err := reader(conn)
	if err != nil {
		return nil, err
	}
	return unpacker(rawRequest, conn)
}

/* ReadRequestFromConn performs the actual reading from
a single connection */
func ReadRequestFromNetConn(conn handle.Conn) ([]byte, error) {
	var packetSize int32
	err := binary.Read(conn, binary.BigEndian, &packetSize)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("prefix read error in readRequestFromConn(): %v", err)
	}
	requestPacket := make([]byte, packetSize)
	_, err = io.ReadFull(conn, requestPacket)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("packet read error in readRequestFromConn(): %v", err)
	}
	return requestPacket, nil
}
