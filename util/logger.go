package util

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type Logger struct {
	name      string
	outLogger *log.Logger
	errLogger *log.Logger
}

const (
	DEBUG = iota
	INFO
	WARN
	ERROR
	FATAL
)

var threadHold = INFO

func init() {
	presets := map[string]int{
		"DEBUG": DEBUG,
		"INFO":  INFO,
		"WARN":  WARN,
		"ERROR": ERROR,
		"FATAL": FATAL,
	}

	level, exist := presets[strings.ToUpper(os.Getenv("LOG_LEVEL"))]

	if exist {
		threadHold = level
	}
}

func (l *Logger) init() *Logger {
	return l
}
func (l *Logger) Debf(str string, v ...interface{}) {
	if threadHold <= DEBUG {
		l.outLogger.Output(2, l.getTime()+" "+fmt.Sprintf(str, v...))
	}

	return
}
func (l *Logger) Logf(str string, v ...interface{}) {
	if threadHold <= INFO {
		l.outLogger.Output(2, l.getTime()+" "+fmt.Sprintf(str, v...))
	}

	return
}
func (l *Logger) Warnf(str string, v ...interface{}) {
	if threadHold <= WARN {
		l.errLogger.Output(2, l.getTime()+" "+fmt.Sprintf(str, v...))
	}

	return
}
func (l *Logger) Errf(str string, v ...interface{}) {
	if threadHold <= INFO {
		l.errLogger.Output(2, l.getTime()+" "+fmt.Sprintf(str, v...))
	}

	return
}
func (l *Logger) Fatalf(str string, v ...interface{}) {
	if threadHold <= FATAL {
		l.errLogger.Output(2, l.getTime()+" "+fmt.Sprintf(str, v...))
	}

	os.Exit(-1)

	return
}

func (l *Logger) getTime() string {
	return time.Now().Format("[2006-01-02 15:04:05]")
}

func NewLogger(component string) (ret *Logger) {
	component = "[" + component + "] "

	ret = &Logger{
		name:      component,
		outLogger: log.New(os.Stdout, component, 0),
		errLogger: log.New(os.Stdout, component, log.Llongfile),
	}

	ret.init()

	return
}
