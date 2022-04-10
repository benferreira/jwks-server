package app_test

import (
	"context"
	"jwks-server/internal/app"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRunTestMode(t *testing.T) {
	defer unsetVars()
	os.Setenv("DEBUG", "true")
	os.Setenv("PORT", "45567")
	os.Setenv("PRETTY_LOGGING", "true")
	os.Setenv("TEST_MODE", "true")

	run(t)
}

func TestRunPubKey(t *testing.T) {
	defer unsetVars()
	os.Setenv("DEBUG", "true")
	os.Setenv("PORT", "45567")
	os.Setenv("RSA_PUB_KEY", `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA0j69to/Sy6u5HeZfBdBt
esvO618W80K3HEo/mnapyhEJsqTQzLy5F+0OM1ZkZCrdFpXuN6jRNHP4tAVEdwg/
aY1uezwbHFwC80HLVNPTmshn5Q4zwFTX60lzSR43BkzwGQPNGixIbkNpLyPBJb0F
yRn9bSbNbwR8GtXYLBNCj4k5ITCb2fWezcxFSajcXfHZPKQgKkXa80P3YEgL/dFU
KFV6gQSmDBm/5nlU4KmUW8loo9AVAL91jW6N2msDlrVg6t8y0aG1w0kAQDwr/sFp
C6W+e6LzeskeuNAvbtgl3blKakSsZ/BB6JIZz+r6YJGFeU0KGYTcpyxSC+T/nJzd
uQIDAQAB
-----END PUBLIC KEY-----`)

	run(t)
}

func run(t *testing.T) {
	application := app.Init()

	go func() {
		client := http.Client{Timeout: time.Duration(1) * time.Second}

		resp, err := client.Get("http://127.0.0.1:45567/api/v1/jwks.json")
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		resp, err = client.Get("http://127.0.0.1:45567/health")
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		application.Server.Shutdown(context.Background())
	}()

	application.Run()
}

func unsetVars() {
	os.Unsetenv("DEBUG")
	os.Unsetenv("PORT")
	os.Unsetenv("PRETTY_LOGGING")
	os.Unsetenv("TEST_MODE")
	os.Unsetenv("RSA_PUB_KEY")
	os.Unsetenv("RSA_KEYS_FILE")
}