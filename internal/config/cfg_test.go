package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_NewConfig(t *testing.T) {
	t.Parallel() // don't do this since we're modifying the environment

	// in case we're called in a batch that sets this value; unsetting
	// happens in a child process, so the batch value is still set for
	// subsequent children
	os.Unsetenv("AUTHN_PATH") // force the error
	cfg, err := NewConfig()
	require.Equal(t, fmt.Errorf("required key AUTHN_PATH missing value"), err)
	require.Nil(t, cfg)

	os.Unsetenv("AUTHN_HOST")         // force AUTHN_PATH to be erased after required:true
	os.Setenv("AUTHN_PATH", "/bogus") // required
	os.Setenv("PHOTO_ROOT", "path with no separator")
	cfg, err = NewConfig()
	require.Nil(t, err)
	require.Equal(t, fmt.Sprintf("path with no separator%s", string(os.PathSeparator)), cfg.PhotoDir)
	require.Empty(t, cfg.AuthnPath, "no host name set, this should be empty")

	os.Setenv("AUTHN_HOST", "bogus") // all is well again
	os.Setenv("PHOTO_ROOT", "normal/")
	os.Setenv("AUTHN_HOST", "localhost") // so AUTHN_PATH doesn't get erased
	cfg, err = NewConfig()
	require.Nil(t, err)
	require.Equal(t, "normal/", cfg.PhotoDir, "making sure separator only gets added once")
	require.Equal(t, "/bogus", cfg.AuthnPath, "now there's a host")
}
