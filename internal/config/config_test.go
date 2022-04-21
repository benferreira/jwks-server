package config_test

import (
	"log"
	"os"
	"testing"

	test_helper "github.com/benferreira/jwks-server/_test_helper"
	"github.com/benferreira/jwks-server/internal/config"

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

	os.Setenv("PORT", "0")
	_, err = config.NewConfigFromEnv()
	assert.NotNil(t, err, "should have errored due to invalid port")
}

func TestNewConfigKeysFromFile(t *testing.T) {
	test_helper.UnsetTestEnvironment()
	defer test_helper.UnsetTestEnvironment()

	path := "../../_test_helper/keys.yaml"
	createKeysFile(path)
	defer deleteKeysFile(path)

	os.Setenv("TEST_MODE", "true")
	os.Setenv("RSA_KEYS_FILE", path)

	config, err := config.NewConfigFromEnv()
	assert.Nil(t, err)
	assert.Equal(t, 2, len(config.PublicKeys.RSAPubKeys))

	assert.NotNil(t, config.PublicKeys.RSAPubKeys[0].Key)
	assert.Equal(t, "key1", config.PublicKeys.RSAPubKeys[0].Kid)

	assert.NotNil(t, config.PublicKeys.RSAPubKeys[1].Key)
	assert.Equal(t, "key2", config.PublicKeys.RSAPubKeys[1].Kid)
}

func TestNewConfigKeysFromFileBad(t *testing.T) {
	test_helper.UnsetTestEnvironment()
	defer test_helper.UnsetTestEnvironment()

	os.Setenv("TEST_MODE", "true")
	os.Setenv("RSA_KEYS_FILE", "./junk/keys.yaml")

	_, err := config.NewConfigFromEnv()
	assert.NotNil(t, err)
}

func TestNewConfigBadTLS(t *testing.T) {
	test_helper.UnsetTestEnvironment()
	defer test_helper.UnsetTestEnvironment()

	os.Setenv("TEST_MODE", "true")
	os.Setenv("TLS", "true")
	os.Setenv("TLS_PRIVATE_KEY_PATH", "./junk/key.pem")
	os.Setenv("TLS_PRIVATE_KEY_PATH", "./junk/cert.pem")

	_, err := config.NewConfigFromEnv()
	assert.NotNil(t, err)
}

func createKeysFile(path string) {
	file, err := os.Create(path)

	if err != nil {
		log.Fatalf("failed to create file %s: %v", path, err)
	}

	file.WriteString(keys)
}

func deleteKeysFile(path string) {
	if err := os.Remove(path); err != nil {
		log.Fatalf("failed to delete file %s, %v", path, err)
	}
}

var keys = `
rsaPubKeys:
  - key: |
      -----BEGIN PUBLIC KEY-----
      MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA0j69to/Sy6u5HeZfBdBt
      esvO618W80K3HEo/mnapyhEJsqTQzLy5F+0OM1ZkZCrdFpXuN6jRNHP4tAVEdwg/
      aY1uezwbHFwC80HLVNPTmshn5Q4zwFTX60lzSR43BkzwGQPNGixIbkNpLyPBJb0F
      yRn9bSbNbwR8GtXYLBNCj4k5ITCb2fWezcxFSajcXfHZPKQgKkXa80P3YEgL/dFU
      KFV6gQSmDBm/5nlU4KmUW8loo9AVAL91jW6N2msDlrVg6t8y0aG1w0kAQDwr/sFp
      C6W+e6LzeskeuNAvbtgl3blKakSsZ/BB6JIZz+r6YJGFeU0KGYTcpyxSC+T/nJzd
      uQIDAQAB
      -----END PUBLIC KEY-----
    kid: key1
  - key: |
      -----BEGIN PUBLIC KEY-----
      MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA9/YhUFEO8X7VkWEbduO8
      gVYS0sHW2G1hyQxVLPY/RPM/mbZ1az3WZTRjRpYtnjQrKn1KHJ671xzn8GfSWRgD
      rkaNF6mTzFI5BTaGORQS+8emFV2xXgVo0pagRFlBdhB2P3xkzNkpJzeKaDD+lG1q
      t1eyJcZTvurSBF9C6nytqdbIc7YIbd8dcYLt3zoCGNm+KBl6aXnpeD6AGt7C1pIj
      yvSWPaY4MUJ6beC6A70PIQxKLgECY9eWLwBqnceUO7EPo9aZuq7tDwE7h7h5TRtT
      pJZ9RTyK6UEiZ/v1e4cKVEVvLL9rmFXLTeJKXYIZbI1lYbuanipiFI/VNvv6r2lF
      CwIDAQAB
      -----END PUBLIC KEY-----
    kid: key2
`
