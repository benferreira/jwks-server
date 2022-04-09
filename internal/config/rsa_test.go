package config_test

import (
	"bytes"
	"jwks-server/internal/config"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRSAPubKeys(t *testing.T) {
	key := "somepublickey"
	keys, err := config.NewRSAPubKeys(key)
	assert.Nil(t, err, "should not err")
	assert.NotNil(t, keys, "should return keys")
	assert.Equal(t, key, keys.RSAPubKeys[0].Key)
}

func TestNewRSAPubKeysInvalidValue(t *testing.T) {
	keys, err := config.NewRSAPubKeys("")
	assert.NotNil(t, err, "should error due to blank public key")
	assert.Nil(t, keys, "should not return keys")
}

func TestNewRSAPubKeysFromFile(t *testing.T) {
	keys, err := config.NewRSAPubKeysFromFile(bytes.NewReader([]byte(valid_rsa_file_text)))
	assert.Nil(t, err, "should not have errored reading contents")
	assert.NotNil(t, keys)
	assert.Equal(t, "key1", keys.RSAPubKeys[0].Kid)
	assert.Equal(t, "key2", keys.RSAPubKeys[1].Kid)
	assert.True(t, strings.Contains(keys.RSAPubKeys[0].Key, "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA0j69to"))
	assert.True(t, strings.Contains(keys.RSAPubKeys[1].Key, "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA9"))
}

func TestNewRSAPubKeysFromFileInvalidKey(t *testing.T) {
	keys, err := config.NewRSAPubKeysFromFile(bytes.NewReader([]byte(invalidkey_rsa_file_text)))
	assert.NotNil(t, err, "should have errored reading contents")
	assert.Nil(t, keys)
}

func TestNewRSAPubKeysFromFileInvalidYaml(t *testing.T) {
	keys, err := config.NewRSAPubKeysFromFile(bytes.NewReader([]byte(invalidyaml_rsa_file_text)))
	assert.NotNil(t, err, "should have errored reading contents")
	assert.Nil(t, keys)
}

var valid_rsa_file_text = `
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

var invalidkey_rsa_file_text = `
rsaPubKeys:
  - kid: key1
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

var invalidyaml_rsa_file_text = `
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
  - bad
`
