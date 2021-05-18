package socketlogger

import (
	"bufio"
	"encoding/json"
	"log"
	"net"
)

type UdpLoggerServer struct {
	LoggerServer
	conn *net.UDPConn
}

func (u *UdpLoggerServer) handleConnections() {
	// Nothing to do here. Udp will listen for all messages and print to the console
}

func (u *UdpLoggerServer) connectSocket(ip string, port int) error {
	sock, err := net.ListenUDP(udpProtocol, &net.UDPAddr{
		Port: port,
		IP:   net.ParseIP(ip),
	})

	if err == nil {
		u.conn = sock
		u.connectionAddr = sock.LocalAddr().String()
	}
	return err
}

// connection coming in will always be nil for the UDP socket, this is applicable to TCP server
func (u *UdpLoggerServer) listenForMsgsOnSocket(net.Conn) {
	reader := bufio.NewReaderSize(u.conn, bufSize)
	dec := json.NewDecoder(reader)
	for {
		msg := SocketMessage{}
		err := dec.Decode(&msg)
		if err != nil {
			er := newMessage(MessageLevelErr, "TcpLoggerServer", "listenForMsgsOnSocket", "ERROR!! %v", err)
			log.Println(er.String())
			break
		} else {
			u.msgs <- msg
		}
	}
}

func NewUdpLoggerServer() loggerServer {
	u := &UdpLoggerServer{}
	u.init(u)
	return u
}
