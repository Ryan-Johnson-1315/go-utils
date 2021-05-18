package socketlogger

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
)

type TcpLoggerServer struct {
	LoggerServer
	serverAddr Connection
}

func NewTcpLoggerServer() loggerServer {
	t := &TcpLoggerServer{}
	t.init(t)
	return t
}

func (t *TcpLoggerServer) handleConnections() {
	l, err := net.Listen("tcp", fmt.Sprintf("%v:%d", t.serverAddr.Addr, t.serverAddr.Port))
	if err != nil {
		log.Fatal(err)
	}

	defer l.Close()
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			er := newMessage(MessageLevelErr, "TcpLoggerServer", "listenForMsgsOnSocket", "Error accepting: %v", err.Error())
			log.Println(er.String())
			continue
		}

		go t.listenForMsgsOnSocket(conn)
	}
}

func (t *TcpLoggerServer) connectSocket(ip string, port int) error {
	t.serverAddr.Addr = ip
	t.serverAddr.Port = port
	return nil // Connecting actually happens in handleConnections
}

func (t *TcpLoggerServer) listenForMsgsOnSocket(conn net.Conn) {
	if conn != nil {
		var sock *net.TCPConn = conn.(*net.TCPConn)
		reader := bufio.NewReaderSize(sock, bufSize)
		dec := json.NewDecoder(reader)
		for {
			msg := SocketMessage{}
			err := dec.Decode(&msg)
			if err != nil {
				er := newMessage(MessageLevelErr, "TcpLoggerServer", "listenForMsgsOnSocket", "ERROR!! %v", err)
				log.Println(er.String())
				sock.Close()
				break
			} else {
				t.msgs <- msg
			}
		}
	}
}
