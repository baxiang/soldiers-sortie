package util

import (
	"bytes"
	"fmt"

	"github.com/sirupsen/logrus"
)

type LogFormatter struct{}

func (l LogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var buffer *bytes.Buffer
	if entry.Buffer != nil {
		buffer = entry.Buffer
	} else {
		buffer = &bytes.Buffer{}
	}
	buffer.Write([]byte("["))
	buffer.Write([]byte(entry.Time.Format("2006-01-02 15:04:05.000")))
	buffer.Write([]byte("] "))
	buffer.Write([]byte("["))
	buffer.Write([]byte(entry.Level.String()))
	buffer.Write([]byte("] "))
	if entry.HasCaller() {
		buffer.Write([]byte("["))
		buffer.Write([]byte(fmt.Sprintf("%s:%d", entry.Caller.File, entry.Caller.Line)))
		buffer.Write([]byte("] "))
	}
	buffer.Write([]byte(entry.Message))
	buffer.Write([]byte("\n"))
	return buffer.Bytes(), nil
}

type ConsulLogger struct{}

func (l *ConsulLogger) Log(args ...interface{}) error {
	logrus.Infoln(args...)
	return nil
}
