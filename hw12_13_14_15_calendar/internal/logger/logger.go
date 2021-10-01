package logger

import (
	"os"
)

const (
	ErrorLevel          = "ERROR"
	WarnLevel           = "WARN"
	InfoLevel           = "INFO"
	DebugLevel          = "DEBUG"
	InfoLogBeginString  = "Info Log: "
	WarnLogBeginString  = "Warn Log: "
	ErrorLogBeginString = "Error Log: "
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
		panic("log file is not found")
	}

	return &Logger{
		Level:   level,
		LogFile: logfile,
		LogLevelMap: LogLevelMap{
			ErrorLevel: 1,
			WarnLevel:  2,
			InfoLevel:  3,
			DebugLevel: 4,
		},
	}
}

func (l Logger) Info(msg string) {
	if l.shouldReportLogLevel(InfoLevel) {
		l.LogFile.WriteString(InfoLogBeginString + msg + "\n")
	}
}

func (l Logger) Warn(msg string) {
	if l.shouldReportLogLevel(WarnLevel) {
		l.LogFile.WriteString(WarnLogBeginString + msg + "\n")
	}
}

func (l Logger) Error(msg string) {
	if l.shouldReportLogLevel(ErrorLevel) {
		l.LogFile.WriteString(ErrorLogBeginString + msg + "\n")
	}
}

func (l Logger) shouldReportLogLevel(requiredLevel string) bool {
	return l.LogLevelMap[l.Level]-l.LogLevelMap[requiredLevel] >= 0
}
