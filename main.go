package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strings"
	"time"

	"github.com/Ryan-Johnson-1315/go-utils/fps"
	"github.com/Ryan-Johnson-1315/go-utils/sockLog"
)

// Utility file to test all of the packages
func main() {
	switch test := strings.Join(os.Args[1:], " "); test {
	case "fps":
		fr := fps.NewFrameRaterWithDescription("Hello world")
		for {
			fr.Tick()
			time.Sleep(15 * time.Millisecond)
		}
	case "logger udp":
		ip, port := "127.0.0.1", 50000
		remoteAddr := &net.UDPAddr{
			Port: port,
			IP:   net.ParseIP(ip),
		}
		if err := sockLog.StartLoggerUDP(ip, port); err == nil {
			sockLog.SetLogFormat(log.Ldate | log.Ltime | log.Lmicroseconds)
			l, err := sockLog.NewUDPSocketLoggerConnection("127.0.0.1", remoteAddr)
			if err != nil {
				log.Fatal(err)
			}
			rand.Seed(time.Now().UnixNano())

			for c := 0; c < 1500; c++ {
				sl := rand.Intn(500)

				msg := sockLog.SocketMessage{
					Caller:      "Test Object",
					MessageType: sockLog.MessageLevel(c % 5),
					Message:     "This is a test message!",
					Function:    "main()",
				}
				l.SendSocketMessage(msg)
				time.Sleep(time.Millisecond * time.Duration(sl))
			}
		}
	case "logger tcp":
		sockLog.StartLoggerTCP(48000)

		service := "127.0.0.1:48000"
		tcpAddr, _ := net.ResolveTCPAddr("tcp", service)

		conn, _ := net.DialTCP("tcp", nil, tcpAddr)

		for c := 0; c < 15; c++ {
			msg := sockLog.SocketMessage{
				Caller:      "Test Object",
				MessageType: sockLog.MessageLevel(c % 5),
				Message:     "This is a test message!",
				Function:    "main()",
			}
			bts, _ := json.Marshal(msg)
			conn.Write(bts)
			time.Sleep(time.Second * 1)
		}
		conn.Close()
		time.Sleep(time.Second * 15)
	default:
		fmt.Println("USAGE")
		fmt.Println("\t go run . [fps | logger udp | logger tcp]")
	}
}
