package logLevel

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

const (
	UNKNOWN = 0
	DEBUG = 1
	INFO = 2
	WARNING = 3
	ERROR = 4
)
var logr = logrus.New()

func init() {
	// Файл для логирования
	log_output, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	if err != nil {
		fmt.Print("Failed to open log file")
		panic(err)
	}

	// Настройки логирования
	logr.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05" ,
	})
	logr.SetLevel(logrus.DebugLevel)
	logr.SetOutput(log_output)
}

func Print(level int, msg string, args ...interface{}) {
	msgf := fmt.Sprintf(msg, args...)
	fmt.Printf("%s: %s;\n", logLvlToStr(level), msgf)

	switch level {
		case DEBUG:
			logr.Debug(msgf)
		case INFO:
			logr.Info(msgf)
		case WARNING:
			logr.Warn(msgf)
		case ERROR:
			logr.Error(msgf)
	}
}

func logLvlToStr(level int) string {
	switch level {
		case DEBUG:
			return "DEBUG"
		case INFO:
			return "INFO"
		case WARNING:
			return "WARNING"
		case ERROR:
			return "ERROR"
		default:
			return "UNKNOWN"
	}
}