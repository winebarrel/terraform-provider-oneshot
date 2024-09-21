package util_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/winebarrel/terraform-provider-oneshot/internal/util"
)

func TestExecCmd_OK(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	stdout, stderr, err := util.ExecCmd("/bin/sh -c", "echo stdout ; echo stderr 1>&2")

	require.NoError(err)
	assert.Equal("stdout\n", stdout)
	assert.Equal("stderr\n", stderr)
}

func TestExecCmd_WithEnv(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	stdout, stderr, err := util.ExecCmd("/bin/sh -c", "echo $FOO ; echo $ZOO 1>&2", "FOO=BAR", "ZOO=BAZ")

	require.NoError(err)
	assert.Equal("BAR\n", stdout)
	assert.Equal("BAZ\n", stderr)
}

func TestExecCmd_Err(t *testing.T) {
	assert := assert.New(t)
	_, _, err := util.ExecCmd("/bin/sh -c", "echo stdout ; echo stderr 1>&2 ; false")
	assert.ErrorContains(err, "Failed to execute command: exit status 1\n[STDOUT] stdout\n\n[STDERR] stderr\n\n")
}
