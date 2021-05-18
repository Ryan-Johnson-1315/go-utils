package socketlogger

import (
	"fmt"
	"log"
	"net"
	"time"
)

type TcpLoggerClient struct {
	socketLogger
	sock *net.TCPConn
}

func (t *TcpLoggerClient) ConnectoToServer(localIp string, localPort int, c Connection) error {
	t.ip = c
	tcpAddr, err := net.ResolveTCPAddr(tcpProtocol, fmt.Sprintf("%s:%d", t.ip.Addr, t.ip.Port))
	if err != nil {
		log.Fatal(err)
	}

	sock, err := net.DialTCP(tcpProtocol, nil, tcpAddr)
	if err != nil {
		log.Fatal(err)
	}

	t.sock = sock
	t.connected = true
	t.Success("TcpLoggerClient", "ConnectToServer", "New TcpLogger Connection -> %s", t.sock.LocalAddr().String())
	time.Sleep(100 * time.Millisecond) // Let everything connect
	return err
}

func NewTcpLoggerClient() Logger {
	a := &TcpLoggerClient{}
	a.init(a)
	return a
}

func (t *TcpLoggerClient) writeOverConnection(msg []byte) error {
	_, err := t.sock.Write(msg)
	if err != nil {
		log.Println(err)
	}
	return err
}
