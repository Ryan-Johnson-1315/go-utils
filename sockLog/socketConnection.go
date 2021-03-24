package sockLog

import (
	"encoding/json"
	"log"
	"net"
	"time"
)

var messageChanneler chan SocketMessage
var disconnected chan bool
var conn *net.TCPConn

type SocketLoggerConnection struct {
}

func NewSocketLoggerConnection(ip string, port int) (SocketLoggerConnection, error) {
	service := "127.0.0.1:48000"
	tcpAddr, _ := net.ResolveTCPAddr("tcp", service)
	var err error
	conn, err = net.DialTCP("tcp", nil, tcpAddr)

	if err != nil {
		log.Fatal("ERROR!! ", err)
		return SocketLoggerConnection{}, err
	}

	messageChanneler = make(chan SocketMessage, 25)
	disconnected = make(chan bool)
	go sendMessages()

	return SocketLoggerConnection{}, nil
}

func (s *SocketLoggerConnection) SendSocketMessage(msg SocketMessage) {
	messageChanneler <- msg
}

func (s *SocketLoggerConnection) SendMessage(caller string, lvl MessageLevel, message, function string) {
	messageChanneler <- SocketMessage{
		Caller:      caller,
		MessageType: lvl,
		Message:     message,
		Function:    function,
	}
}
func sendMessages() {
L:
	for {
		select {
		case msg := <-messageChanneler:
			bts, _ := json.Marshal(msg)
			conn.Write(bts)
			time.Sleep(time.Microsecond * 50) // Flushes buffer
		case <-disconnected:
			break L
		}
	}
	conn.Close()
}
