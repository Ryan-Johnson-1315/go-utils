package socketlogger

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

type loggerServer interface {
	SetLogFlags(flags int)
	SetLogFile(path string)
	Start(ip string, port int) error

	handleConnections()
	connectSocket(ip string, port int) error
	init(i interface{})
	writeMessagesToConsole()
	listenForMsgsOnSocket(net.Conn)
}

type LoggerServer struct {
	Connection
	this           interface{}
	msgs           chan SocketMessage
	connectionAddr string
}

func (l *LoggerServer) SetLogFlags(flags int) {
	log.SetFlags(flags)
}

func (l *LoggerServer) SetLogFile(path string) {
	logFile, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0o666)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(io.MultiWriter(os.Stdout, logFile))
}

func (l *LoggerServer) Start(ip string, port int) error {
	inst, ok := l.this.(loggerServer)
	var err error
	if !ok {
		err = fmt.Errorf("%v is not type 'loggerServer'", l)
	} else {
		err := inst.connectSocket(ip, port)
		if err == nil {
			l.msgs = make(chan SocketMessage, 150)
			go l.writeMessagesToConsole()

			go inst.listenForMsgsOnSocket(nil)
			go inst.handleConnections()
			l.msgs <- newMessage(MessageLevelSuccess, "LoggerServer", "Start()", "Listening for messages @ %s:%d", ip, port)
		}
	}
	return err
}

func (l *LoggerServer) writeMessagesToConsole() {
	for msg := range l.msgs {
		log.Println(msg.String())
	}
}

func (l *LoggerServer) init(i interface{}) {
	l.this = i
}
