package main

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	env := Environment{
		"FOO": {Value: "foo", NeedRemove: false},
		"BAR": {Value: "bar", NeedRemove: false},
	}

	cmd := []string{"/bin/sh", "-c", "echo $FOO $BAR"}

	exitCode := RunCmd(cmd, env)
	require.Equal(t, 0, exitCode)

	// Verify the output
	output, err := exec.Command("/bin/sh", "-c", "echo $FOO $BAR").Output()
	require.NoError(t, err)
	require.Equal(t, "foo bar\n", string(output))
}
