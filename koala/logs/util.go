package logs

import (
	"bytes"
	"fmt"
	"runtime"
	"time"
)

type LogData struct {
	curTime time.Time
	message string
	timeStr string
	level LogLevel
	fileName string
	lineNo int
	traceId string
	serviceName string
	fields *LogField
}

func writeField(buffer *bytes.Buffer,filed, sep string){
	buffer.WriteString(filed)
	buffer.WriteString(sep)
}

func writerKV(buffer *bytes.Buffer,key,val string){
	buffer.WriteString(key)
	buffer.WriteString("=")
	buffer.WriteString(val)
}

func(l *LogData)Bytes()[]byte{
	var buffer bytes.Buffer
	levelStr := getLevelText(l.level)
	writeField(&buffer, l.timeStr, SpaceSep)
	writeField(&buffer, levelStr, SpaceSep)
	writeField(&buffer, l.serviceName, SpaceSep)

	writeField(&buffer, l.fileName, ColonSep)
	writeField(&buffer, fmt.Sprintf("%d", l.lineNo), SpaceSep)
	writeField(&buffer, l.traceId, SpaceSep)

	if l.level == LogLevelAccess && l.fields != nil {
		for _, field := range l.fields.kvs {
			writeField(&buffer, fmt.Sprintf("%v=%v", field.key, field.val), SpaceSep)
		}
	}
	writeField(&buffer, l.message, LineSep)
	return buffer.Bytes()
}

func GetLineInfo() (fileName string, lineNo int) {
	_, fileName, lineNo, _ = runtime.Caller(3)
	return
}