package logger

import (
	"log"
	"os"
)

type LogMsg struct {
	prefix,
	messages []string
}

var (
	loggerChannel = make(chan LogMsg)
	execLog       = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.LUTC|log.Lmicroseconds)
)

func WriteLog(log LogMsg) {
	loggerChannel <- log
}

func ProccessLogs() {

	for msg := range loggerChannel {
		for _, msgLine := range msg.messages {
			execLog.Printf("%s\t%s", msg.prefix, msgLine)
		}
	}
}
