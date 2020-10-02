package route

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"

	"github.com/arstevens/go-request/handle"
)

/* listenAndUnmarshal accepts any connections from listener and attempts
to read and deserialize a request */
func listenAndUnmarshal(listener Listener, unpacker handle.UnpackRequest, reader ReadRequest,
	done <-chan struct{}, outStream chan<- interface{}) {
	defer close(outStream)
	defer listener.Close()

	requestChan := make(chan interface{})
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
attempts to deserialize them into interface{} objects. It then passes
these objects through the returnStream channel */
func receiveRequests(listener Listener, unpacker handle.UnpackRequest, reader ReadRequest, returnStream chan<- interface{}) {
	defer close(returnStream)
	for {
		conn, err := listener.Accept()
		/* An error will be returned if the listener was closed.
		this allows gracefully stopping */
		if err != nil {
			log.Println(err)
			return
		}

		rawRequest, err := reader(conn)
		if err != nil {
			conn.Close()
			log.Println(err)
			continue
		}

		request, err := unpacker(rawRequest)
		if err != nil {
			conn.Close()
			log.Println(err)
			continue
		}
		returnStream <- request
	}
}

/* ReadRequestFromConn performs the actual reading from
a single connection */
func ReadRequestFromNetConn(conn Conn) ([]byte, error) {
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
