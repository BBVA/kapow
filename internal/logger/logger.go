package logger

import (
	"log"
	"os"
)

var L *log.Logger

func init() {
	L = log.New(os.Stderr, "", log.LstdFlags|log.Lmicroseconds|log.LUTC)
}
