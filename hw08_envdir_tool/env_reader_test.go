package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("Test env length", func(t *testing.T) {
		env, err := ReadDir("testdata/env")

		require.NoError(t, err)

		require.Equal(t, 5, len(env))
	})

	t.Run("Directory does not exist", func(t *testing.T) {
		env, err := ReadDir("testdata/env/notexist")

		require.Error(t, err)

		require.Nil(t, env)
	})
}
