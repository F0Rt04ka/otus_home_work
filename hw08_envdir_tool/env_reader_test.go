package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	tmpDir := t.TempDir()

	files := map[string]string{
		"FOO":   "foo",
		"BAR":   " bar \t",
		"EMPTY": "",
		"2LINE": "first\nsecond",
	}

	for name, content := range files {
		err := os.WriteFile(filepath.Join(tmpDir, name), []byte(content), 0o644)
		require.NoError(t, err)
	}

	env, err := ReadDir(tmpDir)
	require.NoError(t, err)

	expected := Environment{
		"FOO":   {Value: "foo", NeedRemove: false},
		"BAR":   {Value: " bar", NeedRemove: false},
		"EMPTY": {Value: "", NeedRemove: true},
		"2LINE": {Value: "first", NeedRemove: false},
	}

	require.Equal(t, expected, env)
}
