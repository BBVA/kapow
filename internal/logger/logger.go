package logger

import (
	"io"
	"log"
	"os"
)

const (
	SCRIPTS = "ScriptsOutput"
)

type LogMsg struct {
	Prefix   string
	Messages []string
}

type internalLogger struct {
	loggerChannel chan LogMsg
	execLog       *log.Logger
}

var loggers = make(map[string]internalLogger)

func RegisterLogger(name string, writer io.Writer) {
	il := internalLogger{}
	flags := log.Ldate | log.Ltime | log.LUTC | log.Lmicroseconds

	if writer == nil {
		writer = os.Stdout
	}

	il.loggerChannel = make(chan LogMsg)
	il.execLog = log.New(writer, "", flags)

	loggers[name] = il
}

func Close(name string) {
	il := loggers[name]

	close(il.loggerChannel)
	il.loggerChannel = nil
}

func SendMsg(name string, log LogMsg) bool {
	if il, ok := loggers[name]; ok {
		il.loggerChannel <- log
		return true
	}

	return false
}

func ProcessMsg(name string) bool {
	var cont bool

	if il, ok := loggers[name]; ok {
		var msg LogMsg

		msg, cont = <-il.loggerChannel

		for _, msgLine := range msg.Messages {
			il.execLog.Printf("%s %s", msg.Prefix, msgLine)
		}
	}

	return cont
}
