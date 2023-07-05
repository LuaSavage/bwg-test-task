package logging

import (
	"fmt"
	"io"
	"os"
	"path"
	"runtime"

	log "github.com/sirupsen/logrus"
)

type writerHook struct {
	Writer    []io.Writer
	LogLevels []log.Level
}

// Levels implements logrus.Hook
func (hook *writerHook) Levels() []log.Level {
	return hook.LogLevels
}

func (hook *writerHook) Fire(entry *log.Entry) error {
	line, err := entry.String()
	if err != nil {
		return err
	}

	for _, w := range hook.Writer {
		_, err = w.Write([]byte(line))
	}
	return err
}

var e *log.Entry

type Logger struct {
	*log.Entry
}

func GetLogger() Logger {
	return Logger{ /*Entry*/ e}
}

// Wraps output message with some custom field=value
func (l *Logger) GetLoggerWithField(k string, v interface{}) Logger {
	return Logger{ /*Entry*/ l.WithField(k, v)}
}

func init() {
	l := log.New()
	l.SetReportCaller( /*reportCaller:*/ true)
	l.Formatter = &log.TextFormatter{
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			filename := path.Base(frame.File)
			return fmt.Sprintf("%s()", frame.Function), fmt.Sprintf("%s:%d", filename, frame.Line)
		},
		DisableColors: false,
		FullTimestamp: true,
	}

	err := os.MkdirAll( /*path*/ "logs" /*permission*/, 0755)
	if err != nil {
		panic(err)
	}

	if err != nil || os.IsExist(err) {
		panic("can't create log dir. no configured logging to files")
	} else {
		logFile, err := os.OpenFile("logs/general.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
		if err != nil {
			panic(fmt.Sprintf("logging: %s", err))
		}

		l.SetOutput(io.Discard)
		l.AddHook(&writerHook{
			Writer:    []io.Writer{logFile, os.Stdout},
			LogLevels: log.AllLevels,
		})
	}

	l.SetLevel(log.TraceLevel)
	e = log.NewEntry(l)
}
