package sockLog

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"
)

// TCPSocketLoggerConnection ...
// Wrapper around the socket that will send both the SocketMessage object
// as well a wrapper for the fields
type TCPSocketLoggerConnection struct {
	messageChanneler chan SocketMessage
	disconnected     chan bool
	conn             *net.TCPConn
}

// NewTCPSocketLoggerConnection ...
// Returns the wrapper or error if something went wrong
func NewTCPSocketLoggerConnection(ip string, port int) (TCPSocketLoggerConnection, error) {
	service := fmt.Sprintf("%s:%d", ip, port)
	tcpAddr, _ := net.ResolveTCPAddr("tcp", service)

	t := TCPSocketLoggerConnection{}

	var err error
	t.conn, err = net.DialTCP("tcp", nil, tcpAddr)

	if err != nil {
		return TCPSocketLoggerConnection{}, err
	}

	t.messageChanneler = make(chan SocketMessage, 25)
	t.disconnected = make(chan bool)
	go t.sendMessages()

	return t, nil
}

// SendSocketMessage ...
// Sends the SocketMessage over the tcp socket
func (t *TCPSocketLoggerConnection) SendSocketMessage(msg SocketMessage) {
	t.messageChanneler <- msg
}

// SendMessage ...
// Builds a SocketMessage and sends it over the tcp socket
func (t *TCPSocketLoggerConnection) SendMessage(caller string, lvl MessageLevel, message, function string) {
	t.messageChanneler <- SocketMessage{
		Caller:      caller,
		MessageType: lvl,
		Message:     message,
		Function:    function,
	}
}

// Close ...
// Closes socket and shuts down goroutines
func (t *TCPSocketLoggerConnection) Close() {
	t.disconnected <- true
}

// UDPSocketLoggerConnection ...
// Wrapper around the socket that will send both the SocketMessage object
// as well a wrapper for the fields
type UDPSocketLoggerConnection struct {
	messageChanneler chan SocketMessage
	disconnected     chan bool
	conn             *net.UDPConn
	remote           *net.UDPAddr
}

// NewUDPSocketLoggerConnection ...
// Returns the wrapper or error if something went wrong
func NewUDPSocketLoggerConnection(ip string, remote *net.UDPAddr) (UDPSocketLoggerConnection, error) {
	var err error
	addr := &net.UDPAddr{
		IP:   net.ParseIP(ip),
		Port: 0,
	}

	u := UDPSocketLoggerConnection{}

	u.conn, err = net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatal("ERROR!! ", err)
		return UDPSocketLoggerConnection{}, err
	}

	u.messageChanneler = make(chan SocketMessage, 25)
	u.disconnected = make(chan bool)
	u.remote = remote
	go u.sendMessages()

	return u, nil
}

// SendSocketMessage ...
// Sends the SocketMessage over the udp socket
func (u *UDPSocketLoggerConnection) SendSocketMessage(msg SocketMessage) {
	u.messageChanneler <- msg
}

// SendMessage ...
// Builds a SocketMessage and sends it over the udp socket
func (u *UDPSocketLoggerConnection) SendMessage(caller string, lvl MessageLevel, message, function string) {
	u.messageChanneler <- SocketMessage{
		Caller:      caller,
		MessageType: lvl,
		Message:     message,
		Function:    function,
	}
}

// Close ...
// Closes socket and shuts down goroutines
func (u *UDPSocketLoggerConnection) Close() {
	u.disconnected <- false
}

//////////////////////////////////////////////////
// Non Exported Utility Functions
//////////////////////////////////////////////////
func (t *TCPSocketLoggerConnection) sendMessages() {
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

func (u *UDPSocketLoggerConnection) sendMessages() {
L:
	for {
		select {
		case msg := <-u.messageChanneler:
			bts, _ := json.Marshal(msg)
			u.conn.WriteToUDP(bts, u.remote)
			time.Sleep(time.Millisecond * 1) // Flushes buffer
		case <-u.disconnected:
			break L
		}
	}
	u.conn.Close()
}
