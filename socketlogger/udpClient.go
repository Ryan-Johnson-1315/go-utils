package socketlogger

import (
	"fmt"
	"log"
	"net"
	"time"
)

type UdpLoggerClient struct {
	socketLogger // Abstract Base class
	conn         *net.UDPConn
	remote       net.Addr
}

func (u *UdpLoggerClient) ConnectoToServer(localIp string, port int, c Connection) error {
	u.ip = c
	remote, err := net.ResolveUDPAddr(udpProtocol, fmt.Sprintf("%s:%d", u.ip.Addr, u.ip.Port))
	if err != nil {
		log.Fatal(err)
	}

	u.remote = remote

	sock, err := net.ListenUDP(udpProtocol, &net.UDPAddr{
		Port: port,
		IP:   net.ParseIP(localIp),
	})
	if err != nil {
		log.Fatal(err)
	}

	u.conn = sock
	u.msgs <- newMessage(MessageLevelSuccess, "UdpClient", "ConnectoToServer", "New UdpLoggerClient Connection -> %v", sock.LocalAddr().String())
	time.Sleep(100 * time.Millisecond) // Let everything connect
	return err
}

func NewUdpLoggerClient() Logger {
	a := &UdpLoggerClient{}
	a.init(a)
	return a
}

func (u *UdpLoggerClient) writeOverConnection(msg []byte) error {
	_, err := u.conn.WriteTo(msg, u.remote)
	if err != nil {
		er := newMessage(MessageLevelErr, "udpLoggerClient", "writeOverConnection", "ERROR!! %v", err)
		log.Println(er.String())
		log.Println(string(msg))
	}
	return err
}
