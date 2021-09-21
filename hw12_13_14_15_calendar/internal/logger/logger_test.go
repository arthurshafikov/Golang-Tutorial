package logger

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	logFilePath := "test.txt"

	t.Run("test debug level", func(t *testing.T) {
		os.Remove(logFilePath)

		logger := New("DEBUG", logFilePath)

		logger.Error("test error")
		logger.Warn("test warn")
		logger.Info("test info")

		expected := "Err log: test error\nWarn log: test warn\nInfo log: test info\n"

		bytes, err := ioutil.ReadFile(logFilePath)
		require.NoError(t, err)
		require.Equal(t, expected, string(bytes))
	})

	t.Run("test error level", func(t *testing.T) {
		os.Remove(logFilePath)

		logger := New("ERROR", logFilePath)

		logger.Error("test error")
		logger.Warn("test warn")
		logger.Info("test info")

		expected := "Err log: test error\n"

		bytes, err := ioutil.ReadFile(logFilePath)
		require.NoError(t, err)
		require.Equal(t, expected, string(bytes))
	})

	t.Run("test warn level", func(t *testing.T) {
		os.Remove(logFilePath)

		logger := New("WARN", logFilePath)

		logger.Error("test error")
		logger.Warn("test warn")
		logger.Info("test info")

		expected := "Err log: test error\nWarn log: test warn\n"

		bytes, err := ioutil.ReadFile(logFilePath)
		require.NoError(t, err)
		require.Equal(t, expected, string(bytes))
	})

	os.Remove(logFilePath)
}
