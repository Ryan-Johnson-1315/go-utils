package main

import (
	"fmt"
	"log"
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
	case "logger tcp":
		sockLog.StartLoggerTCP(48000)

		l1, err := sockLog.NewSocketLoggerConnection("127.0.0.1", 48000)
		if err != nil {
			log.Fatal(err)
		}
		l2, errr := sockLog.NewSocketLoggerConnection("127.0.0.1", 48000)
		if errr != nil {
			log.Fatal(errr)
		}

		for c := 0; c < 15000; c++ {
			l1.SendMessage(sockLog.MessageLevel(c%5), "main.go", "main()", "This is a test %s from l1!!", "message")
			l2.SendMessage(sockLog.MessageLevel(c%5), "main.go", "main()", "This is a test %s from l2!!", "message")
		}
		time.Sleep(time.Second * 15)
	default:
		fmt.Println("USAGE")
		fmt.Println("\t go run . [fps | logger udp | logger tcp]")
	}
}
