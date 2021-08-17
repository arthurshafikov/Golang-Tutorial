package main

import (
	"io/fs"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		envDir := "./testdata/env"
		envMap, err := ReadDir(envDir)

		empty, okEmpty := envMap["EMPTY"]
		bar, okBar := envMap["BAR"]

		require.Len(t, envMap, 5)
		require.NoError(t, err)
		require.True(t, okEmpty)
		require.Equal(t, "", empty.Value)
		require.True(t, okBar)
		require.Equal(t, "bar", bar.Value)
	})

	t.Run("empty env", func(t *testing.T) {
		envDir := "./testdata"
		envMap, err := ReadDir(envDir)

		require.Len(t, envMap, 0)
		require.NoError(t, err)
	})

	t.Run("unexist dir", func(t *testing.T) {
		envDir := "./unexist"
		var pathError *fs.PathError

		envMap, err := ReadDir(envDir)

		require.Len(t, envMap, 0)
		require.ErrorAs(t, err, &pathError)
	})
}
