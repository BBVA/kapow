/*
 * Copyright 2019 Banco Bilbao Vizcaya Argentaria, S.A.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package logger

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"
)

var cleanup = func() { loggers = make(map[string]internalLogger) }

func TestRegisterLoginRegistersWithGivenName(t *testing.T) {
	defer cleanup()
	RegisterLogger("FOO", nil)

	if _, ok := loggers["FOO"]; !ok {
		t.Error("RegisterLogin didn't register the logger")
	}
}

func TestRegisterLoginRegistersWithFlags(t *testing.T) {
	defer cleanup()
	expected := log.Ldate | log.Ltime | log.LUTC | log.Lmicroseconds

	RegisterLogger("FOO", nil)

	if loggers["FOO"].execLog.Flags() != expected {
		t.Errorf("RegisterLogin didn't use correct writer for logger. Expected: %d, got: %d", expected, loggers["FOO"].execLog.Flags())
	}
}

func TestRegisterLoginRegistersWithStdOutWhenNoWriterGiven(t *testing.T) {
	defer cleanup()
	RegisterLogger("FOO", nil)

	if loggers["FOO"].execLog.Writer() != os.Stdout {
		t.Errorf("RegisterLogin didn't use correct writer for logger. Expected: %#v, got: %#v", os.Stdout, loggers["FOO"].execLog.Writer())
	}
}

func TestRegisterLoginRegistersWithWriterGiven(t *testing.T) {
	defer cleanup()
	RegisterLogger("FOO", os.Stderr)

	if loggers["FOO"].execLog.Writer() != os.Stderr {
		t.Errorf("RegisterLogin didn't use correct writer for logger. Expected: %#v, got: %#v", os.Stderr, loggers["FOO"].execLog.Writer())
	}
}

func TestRegisterLoginRegistersWithLogMsgChannel(t *testing.T) {
	defer cleanup()
	expected := reflect.ValueOf(make(chan LogMsg)).String()
	RegisterLogger("FOO", nil)

	if chanType := reflect.ValueOf(loggers["FOO"].loggerChannel).String(); chanType != expected {
		t.Errorf("RegisterLogin didn't create correct channel. Expected: %#v, got: %#v", expected, chanType)
	}
}

func TestCloseClosesLoggerChannel(t *testing.T) {
	defer cleanup()
	RegisterLogger("FOO", nil)
	Close("FOO")

	if _, ok := <-loggers["FOO"].loggerChannel; ok {
		t.Errorf("Close didn't close the channel.")
	}
}

func TestCloseClosesNilsLoggerChannel(t *testing.T) {
	t.Skip("Have to check why it is failing")
	defer cleanup()
	RegisterLogger("FOO", nil)
	Close("FOO")
	fmt.Printf("Channel for logger FOO: %#v\n", loggers["FOO"].loggerChannel)
	if loggers["FOO"].loggerChannel != nil {
		t.Errorf("Close didn't nil the channel.")
	}
}

func TestSendMessageReturnsFalseIfNoLoggerExists(t *testing.T) {
	defer cleanup()

	if ok := SendMsg("FOO", LogMsg{}); ok {
		t.Errorf("SendMessage didn't return error.")
	}
}

func TestSendMessageReturnsTrueAndSendsIfLoggerExists(t *testing.T) {
	t.Skip("Review: there is a DATA RACE error in unit tests")
	defer cleanup()
	var (
		received LogMsg
		msg      LogMsg = LogMsg{"hello", nil}
	)
	RegisterLogger("FOO", nil)
	go func() {
		received = <-loggers["FOO"].loggerChannel
	}()

	ok := SendMsg("FOO", msg)

	if !ok || received.Prefix != msg.Prefix {
		t.Errorf("SendMessage didn't send.")
	}
}

func TestProcessMsgReturnsTrueAfterReceive(t *testing.T) {
	defer cleanup()
	w := &bytes.Buffer{}
	RegisterLogger("FOO", w)
	go SendMsg("FOO", LogMsg{})

	ok := ProcessMsg("FOO")

	if !ok {
		t.Error("ProcessMsg didn't return true")
	}
}

func TestProcessMsgWritesMsgToLog(t *testing.T) {
	expected := []string{"FOOprefix FOO\n", "FOOprefix BAR\n", "FOOprefix FOO BAZ\n"}
	defer cleanup()
	w := &bytes.Buffer{}
	RegisterLogger("FOO", w)
	go SendMsg("FOO", LogMsg{"FOOprefix", []string{"FOO", "BAR", "FOO BAZ"}})

	ProcessMsg("FOO")

	received := w.String()
	for _, ex := range expected {
		if !strings.Contains(received, ex) {
			t.Errorf("ProcessMsg didn't send expected message. Expected: %#v, got: %q", expected, received)
		}
	}
}
