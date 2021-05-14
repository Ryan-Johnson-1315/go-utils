package sockLog

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
)

// SocketLoggerConnection ...
// Wrapper around the socket that will send both the SocketMessage object
// as well a wrapper for the fields
type SocketLoggerConnection struct {
	messageChanneler chan SocketMessage
	disconnected     chan bool
	conn             *net.TCPConn
}

// NewSocketLoggerConnection ...
// Returns the wrapper or error if something went wrong
func NewSocketLoggerConnection(ip string, port int) (SocketLoggerConnection, error) {
	service := fmt.Sprintf("%s:%d", ip, port)
	tcpAddr, _ := net.ResolveTCPAddr("tcp", service)

	t := SocketLoggerConnection{}

	var err error
	t.conn, err = net.DialTCP("tcp", nil, tcpAddr)

	if err != nil {
		return SocketLoggerConnection{}, err
	}

	t.messageChanneler = make(chan SocketMessage, 25)
	t.disconnected = make(chan bool)
	go t.sendMessages()

	return t, nil
}

// SendMessage ...
// Builds a SocketMessage and sends it over the tcp socket
func (t *SocketLoggerConnection) SendMessage(lvl MessageLevel, caller, function, msg string, args ...interface{}) {
	t.messageChanneler <- SocketMessage{
		Caller:      caller,
		Function:    function,
		MessageType: lvl,
		Message:     fmt.Sprintf(msg, args...),
	}
}

// Close ...
// Closes socket and shuts down goroutines
func (t *SocketLoggerConnection) Close() {
	t.disconnected <- true
}

//////////////////////////////////////////////////
// Non Exported Utility Functions
//////////////////////////////////////////////////
func (t *SocketLoggerConnection) sendMessages() {
L:
	for {
		select {
		case msg := <-t.messageChanneler:
			bts, _ := json.Marshal(msg)
			t.conn.Write(bts)
			time.Sleep(time.Millisecond * 1) // Flushes buffer
		case <-t.disconnected:
			break L
		}
	}
	t.conn.Close()
}
