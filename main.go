package main

import (
	"fmt"
	"os"
	"time"

	"github.com/Ryan-Johnson-1315/go-utils/fps"
)

// Utility file to test all of the packages
func main() {
	switch test := os.Args[1]; test {
	case "fps":
		fr := fps.NewFrameRaterWithDescription("Hellow world")
		for {
			fr.Tick()
			time.Sleep(15 * time.Millisecond)
		}
	default:
		fmt.Println("USAGE")
		fmt.Println("\t go run . [fps]")
	}
}
