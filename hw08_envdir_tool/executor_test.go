package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("Panic undefined command name", func(t *testing.T) {
		require.Panics(t, func() {
			emptyArray := make([]string, 0)
			env := Environment{}
			RunCmd(emptyArray, env)
		})
	})
}
