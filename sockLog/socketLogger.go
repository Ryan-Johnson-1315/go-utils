package sockLog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

var logger log.Logger
var messages chan SocketMessage
var logoPrinted bool

func init() {
	logger = log.Logger{}
	log.SetFlags(log.Ldate | log.Ltime)
	messages = make(chan SocketMessage, 500)
	logoPrinted = false
}

// SetLogFormat ...
// Update the log format to a specified format
func SetLogFormat(flags int) {
	log.SetFlags(flags)
}

// StartLogger ...
// Returns error if not started
func StartLoggerUDP(ip string, port int) error {
	sock, err := createSock(ip, port)
	if err != nil {
		return err
	}

	if !logoPrinted {
		fmt.Println(getLogo())
	}

	// UDP Setup
	go printMessages()
	go listenForMessagesUDP(sock)

	log.Printf("%sSocket UDP Logger is starting! Listening @ %s:%d%s", green, ip, port, reset)

	return nil
}

// StartLoggerTCP ...
// This allows the user to know if someone has disconnected
func StartLoggerTCP(port int) {
	go listenForConnections(port)

	if !logoPrinted {
		fmt.Println(getLogo())
		go printMessages()
	}

	log.Printf("%sSocket TCP Logger is starting! Listening @ %s:%d%s", green, "127.0.0.1", port, reset)
	time.Sleep(time.Second * 1) // Wait for socket to connect
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
	remote      string
	connType    string
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

func listenForConnections(port int) {
	l, err := net.Listen("tcp", "localhost:"+fmt.Sprint(port))
	checkErr(err, "", true)

	defer l.Close()
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			continue
		}
		// Handle connections in a new goroutine.<%s>
		go listenForMessagesTCP(conn)
	}
}

func (s *SocketMessage) getColor() string {
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
	return out
}

func (s *SocketMessage) getMessageAsString() string {
	// out := s.getColor()
	if s.remote != "" {
		s.connType = "TCP"
	} else {
		s.connType = "UDP"
	}

	out := fmt.Sprintf("| %s::%s -- %s", s.Caller, s.Function, s.Message)
	return out + reset
}

func listenForMessagesUDP(sock *net.UDPConn) {
	for {
		bts := make([]byte, 32768)
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

func listenForMessagesTCP(conn net.Conn) {
	rmtAddr := conn.RemoteAddr().String()
	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			log.Printf("%sERROR!! Could not read from TCP socket %s:%s%s\n", red, err, rmtAddr, reset)
			conn.Close()
			break
		} else {
			if n > 0 {
				msg := SocketMessage{}
				trimmed := bytes.Trim(buf[0:], "\x00")

				if err := json.Unmarshal(trimmed, &msg); err != nil {
					log.Println("Error unmarshalling:", err)
				} else {
					msg.remote = conn.RemoteAddr().String()
					messages <- msg
				}
			}
		}
	}
	log.Printf("%sERROR!! Closed connection with %s%s\n", red, rmtAddr, reset)
}

func checkErr(err error, description string, shutdown bool) bool {
	if err != nil {
		log.Printf("%sERROR! %s %s%s\n", red, description, err, reset)

		if shutdown {
			os.Exit(1)
		}
	}

	return err == nil
}

func printMessages() {
	for msg := range messages { // Detects if channel is closed
		log.SetPrefix(msg.getColor())
		log.Println(msg.getMessageAsString())
	}
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
	logoPrinted = true
	return logo
}
