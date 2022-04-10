package config_test

import (
	test_helper "jwks-server/_test_helper"
	"jwks-server/internal/config"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	test_helper.UnsetTestEnvironment()
	defer test_helper.UnsetTestEnvironment()

	os.Setenv("DEBUG", "true")
	os.Setenv("PORT", "8080")
	os.Setenv("PRETTY_LOGGING", "true")
	os.Setenv("RSA_PUB_KEY", "somePubKey")

	conf, err := config.NewConfigFromEnv()

	assert.Nil(t, err, "should not error")
	assert.Equal(t, true, conf.Debug)
	assert.Equal(t, 8080, conf.Port)
	assert.Equal(t, true, conf.PrettyLog)
	assert.Equal(t, os.Getenv("RSA_PUB_KEY"), conf.PublicKeys.RSAPubKeys[0].Key)
}

func TestNewConfigMinimumRequiredVars(t *testing.T) {
	test_helper.UnsetTestEnvironment()
	defer test_helper.UnsetTestEnvironment()

	_, err := config.NewConfigFromEnv()
	assert.NotNil(t, err, "should have errored")

	os.Setenv("TEST_MODE", "true")
	conf, err := config.NewConfigFromEnv()
	assert.Nil(t, err, "should not have errored")
	assert.NotNil(t, conf, "config should be returned")
	assert.Equal(t, true, conf.TestMode)

	test_helper.UnsetTestEnvironment()
	os.Setenv("RSA_PUB_KEY", "somePublicKey")
	conf, err = config.NewConfigFromEnv()
	assert.Nil(t, err, "should not have errored")
	assert.NotNil(t, conf, "config should be returned")
	assert.Equal(t, os.Getenv("RSA_PUB_KEY"), conf.PublicKeys.RSAPubKeys[0].Key)
}

func TestNewConfigBadPort(t *testing.T) {
	test_helper.UnsetTestEnvironment()
	defer test_helper.UnsetTestEnvironment()

	os.Setenv("PORT", "badPort")
	_, err := config.NewConfigFromEnv()
	assert.NotNil(t, err, "should have errored due to invalid port")
}
