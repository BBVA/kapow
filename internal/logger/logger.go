package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

var A *log.Logger
var L *log.Logger

func init() {
	A = log.New(os.Stdout, "", 0)
	L = log.New(os.Stderr, "", log.LstdFlags|log.Lmicroseconds|log.LUTC)
}

// See https://en.wikipedia.org/wiki/Common_Log_Format
// Tested with love against CLFParser https://pypi.org/project/clfparser/
func LogAccess(clientAddr, handlerId, user, method, resource, version string, status int, bytesSent int64, referer, userAgent string) {
	clientAddr = dashIfEmpty(clientAddr)
	handlerId = dashIfEmpty(handlerId)
	user = dashIfEmpty(user)
	referer = dashIfEmpty(referer)
	userAgent = dashIfEmpty(userAgent)
	firstLine := fmt.Sprintf("%s %s %s", method, resource, version)
	ts := time.Now().UTC().Format("02/Jan/2006:15:04:05 -0700") // Amazing date format layout!  Not.
	A.Printf("%s %s %s [%s] %q %d %d %q %q\n",
		clientAddr,
		handlerId,
		user,
		ts,
		firstLine,
		status,
		bytesSent,
		referer,
		userAgent)
}

func dashIfEmpty(value string) string {
	if value == "" {
		return "-"
	}
	return value
}
