package sockLog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
)

var logger log.Logger
var messages chan SocketMessage

func init() {
	logger = log.Logger{}
	log.SetFlags(log.Ldate | log.Ltime)
	messages = make(chan SocketMessage, 100)
}

// SetLogFormat ...
// Update the log format to a specified format
func SetLogFormat(flags int) {
	log.SetFlags(flags)
}

// StartLogger ...
// Returns error if not started
func StartLogger(ip string, port int) error {
	sock, err := createSock(ip, port)
	if err != nil {
		return err
	}

	go printMessages()
	go listenForMessages(sock)

	fmt.Println(getLogo())
	log.Printf("%sSocket Logger is starting! Listening @ %s:%d%s", green, ip, port, reset)

	return nil
}

type MessageLevel int

const (
	MessageLevelLog     MessageLevel = 0
	MessageLevelWarn    MessageLevel = 1
	MessageLevelSuccess MessageLevel = 2
	MessageLevelErr     MessageLevel = 3
	MessageLevelDebug   MessageLevel = 4
)

// SocketMessage ...
// Message struct that the JSON will unmarshal to
type SocketMessage struct {
	Caller      string       `json:"caller"`
	MessageType MessageLevel `json:"level"`
	Message     string       `json:"message"`
	Function    string       `json:"function"`
}

//////////////////////////////////////////////////
// Non Exported Utility Functions
//////////////////////////////////////////////////
const (
	reset  string = "\033[0m"
	red    string = "\033[31m"
	green  string = "\033[32m"
	yellow string = "\033[33m"
	cyan   string = "\033[36m"

	udpProtocol string = "udp"
)

func (s *SocketMessage) getMessageAsString() string {
	out := ""
	switch s.MessageType {
	case MessageLevelWarn:
		out += yellow
		break
	case MessageLevelSuccess:
		out += green
		break
	case MessageLevelErr:
		out += red
		break
	case MessageLevelDebug:
		out += cyan
		break
	default: // MessageLevelLog
		out += reset
	}

	out += fmt.Sprintf("<%s::%s>  %s", s.Caller, s.Function, s.Message)
	return out + reset
}

func listenForMessages(sock *net.UDPConn) {
	for {
		bts := make([]byte, 1024)
		sock.ReadFromUDP(bts)
		if len(bts) > 0 {
			trimmed := bytes.Trim(bts, "\x00")
			newMessage := SocketMessage{}

			if err := json.Unmarshal(trimmed, &newMessage); err != nil {
				log.Println("Error unmarshalling:", err)
			} else {
				messages <- newMessage
			}
		}
	}
}

func printMessages() {
	for msg := range messages { // Detects if channel is closed
		log.Println(msg.getMessageAsString())
	}
	log.Println("Clossing down printMessages()")
}

func createIPAddr(ip string, port int) *net.UDPAddr {
	return &net.UDPAddr{
		Port: port,
		IP:   net.ParseIP(ip),
	}
}

func createSock(ip string, port int) (*net.UDPConn, error) {
	return net.ListenUDP(udpProtocol, createIPAddr(ip, port))
}

func getLogo() string {
	logo := cyan
	logo += "///////////////////////////////////////// \\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n"
	logo += `  _________              __           __   ____                                      
 /   _____/ ____   ____ |  | __ _____/  |_|    |    ____   ____   ____   ___________ 
 \_____  \ /  _ \_/ ___\|  |/ // __ \   __\    |   /  _ \ / ___\ / ___\_/ __ \_  __ \
 /        (  <_> )  \___|    <\  ___/|  | |    |__(  <_> ) /_/  > /_/  >  ___/|  | \/
/_______  /\____/ \___  >__|_ \\___  >__| |_______ \____/\___  /\___  / \___  >__|   
        \/            \/     \/    \/             \/    /_____//_____/      \/       `
	logo += "\n\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\ ///////////////////////////////////////////\n"
	logo += reset
	return logo
}
