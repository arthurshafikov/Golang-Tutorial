package logger

import (
	"log"
	"os"
)

type LogLevelMap map[string]int

type Logger struct {
	Level       string
	LogFile     *os.File
	LogLevelMap LogLevelMap
}

func New(level string, logFilePath string) *Logger {
	logfile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		// todo
		log.Println(err)
		panic("log file is not found")
	}
	return &Logger{
		Level:   level,
		LogFile: logfile,
		LogLevelMap: LogLevelMap{
			"ERROR": 1,
			"WARN":  2,
			"INFO":  3,
			"DEBUG": 4,
		},
	}
}

func (l Logger) Info(msg string) {
	if l.checkLogLevel("INFO") {
		l.LogFile.WriteString("Info log: " + msg + "\n")
		// log.Println("Info log: ", msg)
	}
}

func (l Logger) Warn(msg string) {
	if l.checkLogLevel("WARN") {
		l.LogFile.WriteString("Warn log: " + msg + "\n")
		// log.Println("Warn log: ", msg)
	}
}

func (l Logger) Error(msg string) {
	if l.checkLogLevel("ERROR") {
		l.LogFile.WriteString("Err log: " + msg + "\n")
		// log.Println("Err log: ", msg)
	}
}

func (l Logger) checkLogLevel(requiredLevel string) bool {
	return l.LogLevelMap[l.Level]-l.LogLevelMap[requiredLevel] >= 0
}
