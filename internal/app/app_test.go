package app_test

import (
	"context"
	"fmt"
	test_helper "jwks-server/_test_helper"
	"jwks-server/internal/app"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRunTestMode(t *testing.T) {
	test_helper.UnsetTestEnvironment()
	defer test_helper.UnsetTestEnvironment()
	os.Setenv("DEBUG", "true")
	os.Setenv("PORT", "45567")
	os.Setenv("PRETTY_LOGGING", "true")
	os.Setenv("TEST_MODE", "true")

	run(t)
}

func TestRunPubKey(t *testing.T) {
	test_helper.UnsetTestEnvironment()
	defer test_helper.UnsetTestEnvironment()
	os.Setenv("DEBUG", "true")
	os.Setenv("PORT", "45568")
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
		time.Sleep(1 * time.Second)
		client := http.Client{Timeout: time.Duration(1) * time.Second}

		baseUrl := fmt.Sprintf("http://127.0.0.1:%d", application.Configuration.Port)

		resp, err := client.Get(fmt.Sprintf("%s/api/v1/jwks.json", baseUrl))
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		resp, err = client.Get(fmt.Sprintf("%s/health", baseUrl))
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		application.Server.Shutdown(context.Background())
	}()

	application.Run()
}
