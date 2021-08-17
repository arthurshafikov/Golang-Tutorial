package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		cmd := []string{
			"/bin/bash",
			"$(pwd)/testdata/echo.sh",
			"arg1=1 arg2=2",
		}

		err := RunCmd(cmd, nil)

		require.Equal(t, 0, err)
	})
}
