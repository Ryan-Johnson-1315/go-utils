package sockLog

import (
	"encoding/json"
	"log"
	"net"
	"time"
)

var messageChannelerTCP chan SocketMessage
var disconnectedTCP chan bool
var connTCP *net.TCPConn

type TCPSocketLoggerConnection struct {
}

func NewTCPSocketLoggerConnection(ip string, port int) (TCPSocketLoggerConnection, error) {
	service := "127.0.0.1:48000"
	tcpAddr, _ := net.ResolveTCPAddr("tcp", service)
	var err error
	connTCP, err = net.DialTCP("tcp", nil, tcpAddr)

	if err != nil {
		log.Fatal("ERROR!! ", err)
		return TCPSocketLoggerConnection{}, err
	}

	messageChannelerTCP = make(chan SocketMessage, 25)
	disconnectedTCP = make(chan bool)
	go sendMessagesTCP()

	return TCPSocketLoggerConnection{}, nil
}

func (s *TCPSocketLoggerConnection) SendSocketMessage(msg SocketMessage) {
	messageChannelerTCP <- msg
}

func (s *TCPSocketLoggerConnection) SendMessage(caller string, lvl MessageLevel, message, function string) {
	messageChannelerTCP <- SocketMessage{
		Caller:      caller,
		MessageType: lvl,
		Message:     message,
		Function:    function,
	}
}
func sendMessagesTCP() {
L:
	for {
		select {
		case msg := <-messageChannelerTCP:
			bts, _ := json.Marshal(msg)
			connTCP.Write(bts)
			time.Sleep(time.Microsecond * 50) // Flushes buffer
		case <-disconnectedTCP:
			break L
		}
	}
	connTCP.Close()
}
