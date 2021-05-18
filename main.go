package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Ryan-Johnson-1315/go-utils/fps"
	"github.com/Ryan-Johnson-1315/go-utils/socketlogger"
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
	case "logger":
		udp := socketlogger.NewUdpLoggerServer()
		udp.SetLogFile("output.log")
		udp.Start("127.0.0.1", 12234)

		time.Sleep(1 * time.Second)
		udpL := socketlogger.NewUdpLoggerClient()
		udpL.ConnectoToServer("127.0.0.1", 0, socketlogger.Connection{Addr: "127.0.0.1", Port: 12234})
		udpLL := socketlogger.NewUdpLoggerClient()
		udpLL.ConnectoToServer("127.0.0.1", 0, socketlogger.Connection{Addr: "127.0.0.1", Port: 12234})

		tcp := socketlogger.NewTcpLoggerServer()
		tcp.SetLogFile("output.log")
		tcp.Start("127.0.0.1", 12235)

		time.Sleep(1 * time.Second)
		tcpL := socketlogger.NewTcpLoggerClient()
		tcpL.ConnectoToServer("127.0.0.1", 0, socketlogger.Connection{Addr: "127.0.0.1", Port: 12235})

		tcpLL := socketlogger.NewTcpLoggerClient()
		tcpLL.ConnectoToServer("127.0.0.1", 0, socketlogger.Connection{Addr: "127.0.0.1", Port: 12235})

		go func() {
			for i := 0; i < 10; i++ {
				switch i % 5 {
				case 0:
					udpL.Log("main", "main()", "Testing testing! Hello udp %s", "world")
					udpLL.Log("main", "main()", "Testing testing! Hello udp %s", "world")
				case 2:
					udpL.Wrn("main", "main()", "Testing testing! Hello udp %s", "world")
					udpLL.Wrn("main", "main()", "Testing testing! Hello udp %s", "world")
				case 1:
					udpL.Dbg("main", "main()", "Testing testing! Hello udp %s", "world")
					udpLL.Dbg("main", "main()", "Testing testing! Hello udp %s", "world")
				case 4:
					udpL.Success("main", "main()", "Testing testing! Hello udp %s", "world")
					udpLL.Success("main", "main()", "Testing testing! Hello udp %s", "world")
				case 3:
					udpL.Err("main", "main()", "Testing testing! Hello udp %s", "world")
					udpLL.Err("main", "main()", "Testing testing! Hello udp %s", "world")
				}
				time.Sleep(time.Duration(5) * time.Millisecond)
			}
		}()

		go func() {
			for i := 0; i < 10; i++ {
				switch i % 5 {
				case 0:
					tcpL.Log("main", "main()", "Testing testing! Hello tcp %s", "world")
					tcpLL.Log("main", "main()", "Testing testing! Hello tcp %s", "world")
				case 4:
					tcpL.Wrn("main", "main()", "Testing testing! Hello tcp %s", "world")
					tcpLL.Wrn("main", "main()", "Testing testing! Hello tcp %s", "world")
				case 3:
					tcpL.Dbg("main", "main()", "Testing testing! Hello tcp %s", "world")
					tcpLL.Dbg("main", "main()", "Testing testing! Hello tcp %s", "world")
				case 2:
					tcpL.Success("main", "main()", "Testing testing! Hello tcp %s", "world")
					tcpLL.Success("main", "main()", "Testing testing! Hello tcp %s", "world")
				case 1:
					tcpL.Err("main", "main()", "Testing testing! Hello tcp %s", "world")
					tcpLL.Err("main", "main()", "Testing testing! Hello tcp %s", "world")
				}
				time.Sleep(time.Duration(5) * time.Millisecond)
			}
		}()
		time.Sleep(1 * time.Second)
	default:
		fmt.Println("USAGE")
		fmt.Println("\t go run . [fps | logger]")
	}
}
