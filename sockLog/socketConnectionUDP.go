package sockLog

import (
	"encoding/json"
	"log"
	"net"
	"time"
)

var messageChannelerUDP chan SocketMessage
var disconnectedUDP chan bool
var connUDP *net.UDPConn
var remoteUDP *net.UDPAddr

type UDPSocketLoggerConnection struct {
}

func NewUDPSocketLoggerConnection(ip string, remote *net.UDPAddr) (UDPSocketLoggerConnection, error) {
	var err error
	addr := &net.UDPAddr{
		IP:   net.ParseIP(ip),
		Port: 0,
	}
	connUDP, err = net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatal("ERROR!! ", err)
		return UDPSocketLoggerConnection{}, err
	}

	messageChannelerUDP = make(chan SocketMessage, 25)
	disconnectedUDP = make(chan bool)
	remoteUDP = remote
	go sendMessagesUDP()

	return UDPSocketLoggerConnection{}, nil
}

func (s *UDPSocketLoggerConnection) SendSocketMessage(msg SocketMessage) {
	messageChannelerUDP <- msg
}

func (s *UDPSocketLoggerConnection) SendMessage(caller string, lvl MessageLevel, message, function string) {
	messageChannelerUDP <- SocketMessage{
		Caller:      caller,
		MessageType: lvl,
		Message:     message,
		Function:    function,
	}
}
func sendMessagesUDP() {
L:
	for {
		select {
		case msg := <-messageChannelerUDP:
			bts, _ := json.Marshal(msg)
			connUDP.WriteToUDP(bts, remoteUDP)
			time.Sleep(time.Microsecond * 50) // Flushes buffer
		case <-disconnectedUDP:
			break L
		}
	}
	connUDP.Close()
}
