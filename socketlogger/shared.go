package socketlogger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
)

type (
	messageLevel int
	color        string
)

const (
	reset  color = "\033[0m"
	red    color = "\033[31m"
	green  color = "\033[32m"
	yellow color = "\033[33m"
	cyan   color = "\033[36m"

	udpProtocol string = "udp"
	tcpProtocol string = "tcp"

	MessageLevelLog     messageLevel = 0
	MessageLevelWarn    messageLevel = 1
	MessageLevelSuccess messageLevel = 2
	MessageLevelErr     messageLevel = 3
	MessageLevelDebug   messageLevel = 4

	bufSize int = 16384
)

type Connection struct {
	Addr string
	Port int
}

type SocketMessage struct {
	Caller      string       `json:"caller"`
	MessageType messageLevel `json:"level"`
	Message     string       `json:"message"`
	Function    string       `json:"function"`
}

func (s *SocketMessage) String() string {
	str := string(reset)
	switch s.MessageType {
	case MessageLevelWarn:
		str += string(yellow)
	case MessageLevelSuccess:
		str += string(green)
	case MessageLevelErr:
		str += string(red)
	case MessageLevelDebug:
		str += string(cyan)
	}
	str += fmt.Sprintf(" | %s::%s -- %s", s.Caller, s.Function, s.Message)

	return str + string(reset)
}

func (s *SocketMessage) asBytes() ([]byte, error) {
	return json.Marshal(s)
}

func parseFromBytes(bts []byte) []SocketMessage {
	msgs := make([]SocketMessage, 0)
	for _, buf := range bytes.Split(bts, []byte("\n")) {
		if len(buf) > 0 {
			msg := SocketMessage{}
			err := json.Unmarshal(buf, &msg)
			if err != nil {
				log.Println(red, "shared::parseFromBytes --", err, string(buf), reset)
			} else {
				msgs = append(msgs, msg)
			}
		}
	}
	return msgs
}

func newMessage(lvl messageLevel, caller, function, format string, args ...interface{}) SocketMessage {
	return SocketMessage{
		Caller:      caller,
		MessageType: lvl,
		Message:     fmt.Sprintf(format, args...),
		Function:    function,
	}
}
