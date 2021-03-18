package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/Ryan-Johnson-1315/go-utils/fps"
	"github.com/Ryan-Johnson-1315/go-utils/sockLog"
)

// Utility file to test all of the packages
func main() {
	switch test := os.Args[1]; test {
	case "fps":
		fr := fps.NewFrameRaterWithDescription("Hello world")
		for {
			fr.Tick()
			time.Sleep(15 * time.Millisecond)
		}
	case "logger":
		ip, port := "127.0.0.1", 50000
		remoteAddr := &net.UDPAddr{
			Port: port,
			IP:   net.ParseIP(ip),
		}
		if err := sockLog.StartLogger(ip, port); err == nil {
			addr := &net.UDPAddr{
				Port: 50001,
				IP:   net.ParseIP(ip),
			}

			sock, _ := net.ListenUDP("udp", addr)
			count := 0
			for {
				count++
				msg := sockLog.SocketMessage{
					Caller:      "Test Object",
					MessageType: sockLog.MessageLevel(count % 5),
					Message:     "This is a test message!",
					Function:    "main()",
				}
				bts, _ := json.Marshal(msg)
				sock.WriteToUDP(bts, remoteAddr)
				time.Sleep(time.Second * 1)
			}
		}
	default:
		fmt.Println("USAGE")
		fmt.Println("\t go run . [fps]")
	}
}
