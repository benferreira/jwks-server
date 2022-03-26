package config_test

import (
	"jwks-server/internal/config"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	unsetVars()
	defer unsetVars()

	os.Setenv("DEBUG", "true")
	os.Setenv("PORT", "8080")
	os.Setenv("PRETTY_LOGGING", "true")
	os.Setenv("RSA_PUB_KEY", "somePubKey")

	conf, err := config.NewConfig()

	assert.Nil(t, err, "should not error")
	assert.Equal(t, true, conf.Debug)
	assert.Equal(t, 8080, conf.Port)
	assert.Equal(t, true, conf.PrettyLog)
	assert.Equal(t, os.Getenv("RSA_PUB_KEY"), conf.RsaPubKey)
}

func TestNewConfigMinimumRequiredVars(t *testing.T) {
	unsetVars()
	defer unsetVars()

	_, err := config.NewConfig()
	assert.NotNil(t, err, "should have errored")

	os.Setenv("TEST_MODE", "true")
	conf, err := config.NewConfig()
	assert.Nil(t, err, "should not have errored")
	assert.NotNil(t, conf, "config should be returned")
	assert.Equal(t, true, conf.TestMode)

	unsetVars()
	os.Setenv("RSA_PUB_KEY", "somePublicKey")
	conf, err = config.NewConfig()
	assert.Nil(t, err, "should not have errored")
	assert.NotNil(t, conf, "config should be returned")
	assert.Equal(t, os.Getenv("RSA_PUB_KEY"), conf.RsaPubKey)
}

func TestNewConfigBadPort(t *testing.T) {
	unsetVars()
	defer unsetVars()

	os.Setenv("PORT", "badPort")
	_, err := config.NewConfig()
	assert.NotNil(t, err, "should have errored due to invalid port")
}

func unsetVars() {
	os.Unsetenv("DEBUG")
	os.Unsetenv("PORT")
	os.Unsetenv("PRETTY_LOGGING")
	os.Unsetenv("TEST_MODE")
	os.Unsetenv("RSA_PUB_KEY")
}
