package main

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	fromPath = "tmp/in.txt"
	toPath   = "tmp/out.txt"
)

func TestCopy(t *testing.T) {
	displayProgressBar = false
	t.Run("negative arguments", func(t *testing.T) {
		limit = -1
		offset = -1
		err := Copy(fromPath, toPath, offset, limit)

		require.ErrorIs(t, err, ErrWrongArguments)
	})

	t.Run("fromPath does not exist", func(t *testing.T) {
		fromPath := "tmp/not-exist.txt"
		err := Copy(fromPath, toPath, 0, 0)

		require.ErrorIs(t, err, ErrFileDoesNotExists)
	})

	t.Run("offset more than size of file", func(t *testing.T) {
		limit = 0
		offset = 999999999
		err := Copy(fromPath, toPath, offset, limit)

		require.ErrorIs(t, err, ErrOffsetExceedsFileSize)
	})

	t.Run("limit more than size of file", func(t *testing.T) {
		limit = 999999999
		offset = 0
		err := Copy(fromPath, toPath, offset, limit)

		fromFile, _ := ioutil.ReadFile(fromPath)
		toFile, _ := ioutil.ReadFile(toPath)

		require.Equal(t, fromFile, toFile)
		require.NoError(t, err)
	})

	t.Run("zero limit zero offset", func(t *testing.T) {
		limit = 0
		offset = 0
		err := Copy(fromPath, toPath, offset, limit)

		fromFile, _ := ioutil.ReadFile(fromPath)
		toFile, _ := ioutil.ReadFile(toPath)

		require.Equal(t, fromFile, toFile)
		require.NoError(t, err)
	})

	t.Run("five limit ten offset", func(t *testing.T) {
		limit = 5
		offset = 10
		expected := "m dol"
		err := Copy(fromPath, toPath, offset, limit)

		toFile, _ := ioutil.ReadFile(toPath)

		require.Equal(t, expected, string(toFile))
		require.NoError(t, err)
	})
}

func TestProgressBar(t *testing.T) {
	t.Run("process", func(t *testing.T) {
		pb := NewProgressBar(200, false)

		pb.Process(10)
		pb.Process(20)

		require.Equal(t, uint(30), pb.current)
		require.Equal(t, uint(15), pb.findPercent())
	})
}
