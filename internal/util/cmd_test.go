package util_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/winebarrel/terraform-provider-oneshot/internal/util"
)

func TestCmdRun_OK(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	cmd := util.NewCmd("/bin/bash -c", "/dev/null", "/dev/null")
	stdout, stderr, err := cmd.Run("echo stdout ; echo stderr 1>&2")

	require.NoError(err)
	assert.Equal("stdout\n", stdout)
	assert.Equal("stderr\n", stderr)
}

func TestCmdRun_WithEnv(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	cmd := util.NewCmd("/bin/bash -c", "/dev/null", "/dev/null")
	stdout, stderr, err := cmd.Run("echo $FOO ; echo $ZOO 1>&2", "FOO=BAR", "ZOO=BAZ")

	require.NoError(err)
	assert.Equal("BAR\n", stdout)
	assert.Equal("BAZ\n", stderr)
}

func TestCmdRun_Err(t *testing.T) {
	assert := assert.New(t)
	cmd := util.NewCmd("/bin/bash -c", "/dev/null", "/dev/null")
	_, _, err := cmd.Run("echo stdout ; echo stderr 1>&2 ; false")
	assert.ErrorContains(err, "Failed to execute command: exit status 1\n[STDOUT] stdout\n\n[STDERR] stderr\n\n")
}

func TestCmdRun_WithLog(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	cwd, _ := os.Getwd()
	os.Chdir(t.TempDir())
	defer os.Chdir(cwd)

	cmd := util.NewCmd("/bin/bash -c", "stdout.log", "stderr.log")
	stdout, stderr, err := cmd.Run("echo stdout ; echo stderr 1>&2")

	require.NoError(err)
	assert.Equal("stdout\n", stdout)
	assert.Equal("stderr\n", stderr)

	stdoutLog, _ := os.ReadFile("stdout.log")
	assert.Equal("stdout\n", string(stdoutLog))
	stderrLog, _ := os.ReadFile("stderr.log")
	assert.Equal("stderr\n", string(stderrLog))
}
