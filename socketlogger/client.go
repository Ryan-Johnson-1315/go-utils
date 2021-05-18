package socketlogger

import (
	"fmt"
	"log"
	"time"
)

// This file contains the interface and implementation of what is considered
// the 'logger'. This differs from the log server.

type Logger interface {
	ConnectoToServer(localIP string, port int, addr Connection) error
	Success(caller, function, format string, args ...interface{})
	Log(caller, function, format string, args ...interface{})
	Wrn(caller, function, format string, args ...interface{})
	Err(caller, function, format string, args ...interface{})
	Dbg(caller, function, format string, args ...interface{})

	writeOverConnection([]byte) error
	sendFromChannel()
}

type socketLogger struct {
	ip        Connection
	msgs      chan SocketMessage
	this      interface{}
	connected bool
}

func (l *socketLogger) Log(caller, function, format string, args ...interface{}) {
	l.msgs <- newMessage(MessageLevelLog, caller, function, format, args...)
}

func (l *socketLogger) Wrn(caller, function, format string, args ...interface{}) {
	l.msgs <- newMessage(MessageLevelWarn, caller, function, format, args...)
}

func (l *socketLogger) Success(caller, function, format string, args ...interface{}) {
	l.msgs <- newMessage(MessageLevelSuccess, caller, function, format, args...)
}

func (l *socketLogger) Err(caller, function, format string, args ...interface{}) {
	l.msgs <- newMessage(MessageLevelErr, caller, function, format, args...)
}

func (l *socketLogger) Dbg(caller, function, format string, args ...interface{}) {
	l.msgs <- newMessage(MessageLevelDebug, caller, function, format, args...)
}

func (l *socketLogger) init(instance interface{}) {
	l.msgs = make(chan SocketMessage, 150)
	l.connected = false
	l.this = instance
	go l.sendFromChannel() // Different based on TCP/UDP/Serial
}

func (l *socketLogger) sendFromChannel() {
	lg, ok := l.this.(Logger)
	if !ok {
		panic(fmt.Sprintf("%v is not a 'logger' type", lg))
	} else {
		for msg := range l.msgs {
			bytes, err := msg.asBytes()
			if err != nil {
				log.Println("ERROR!!", err)
			} else {
				if err := lg.writeOverConnection(bytes); err != nil {
					log.Println(msg)
				}
				time.Sleep(2 * time.Millisecond)
			}
		}
	}
}
