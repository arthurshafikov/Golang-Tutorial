package logger

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const logFilePath = "test.txt"

func TestDebugLevel(t *testing.T) {
	logger := New(DebugLevel, logFilePath)

	logger.Error("test error")
	logger.Warn("test warn")
	logger.Info("test info")

	bytes, err := ioutil.ReadFile(logFilePath)
	require.NoError(t, err)

	expected := ErrorLogBeginString + "test error\n" + WarnLogBeginString + "test warn\n" + InfoLogBeginString + "test info\n"
	require.Equal(t, expected, string(bytes))

	os.Remove(logFilePath)
}

func TestErrorLevel(t *testing.T) {
	logger := New(ErrorLevel, logFilePath)

	logger.Error("test error")
	logger.Warn("test warn")
	logger.Info("test info")

	bytes, err := ioutil.ReadFile(logFilePath)
	require.NoError(t, err)

	expected := ErrorLogBeginString + "test error\n"
	require.Equal(t, expected, string(bytes))

	os.Remove(logFilePath)
}

func TestWarnLevel(t *testing.T) {
	logger := New(WarnLevel, logFilePath)

	logger.Error("test error")
	logger.Warn("test warn")
	logger.Info("test info")

	bytes, err := ioutil.ReadFile(logFilePath)
	require.NoError(t, err)

	expected := ErrorLogBeginString + "test error\n" + WarnLogBeginString + "test warn\n"
	require.Equal(t, expected, string(bytes))

	os.Remove(logFilePath)
}
