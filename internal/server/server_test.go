package server_test

import (
	"context"
	"io/ioutil"
	"jwks-server/internal/config"
	"jwks-server/internal/server"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestServe(t *testing.T) {
	serv, err := server.NewServer(&config.ServerConfig{Port: 45566}, `{"test":"json"}`)
	assert.Nil(t, err)

	go func() {
		client := http.Client{Timeout: time.Duration(1) * time.Second}

		resp, err := client.Get("http://127.0.0.1:45566/health")
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		serv.Shutdown(context.Background())
	}()

	err = serv.Start()
	assert.Nil(t, err)
}

func TestServePortCollision(t *testing.T) {
	conf := config.ServerConfig{Port: 45566}
	serv, _ := server.NewServer(&conf, `{"test":"json"}`)
	serv2, _ := server.NewServer(&conf, `{"test":"json"}`)

	go func() {
		err := serv2.Start()
		assert.NotNil(t, err, "should error due to port collision")
		serv.Shutdown(context.Background())
		serv2.Shutdown(context.Background())
	}()

	serv.Start()
}

func TestHealth(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	server.HealthCheckHandler(w, r)
	result := w.Result()

	assert.Equal(t, http.StatusOK, result.StatusCode)
}

func TestHealthBadRequest(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/health", nil)
	w := httptest.NewRecorder()

	server.HealthCheckHandler(w, r)
	result := w.Result()

	assert.Equal(t, http.StatusBadRequest, result.StatusCode)
}

func TestJWKS(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/api/v1/jwks.json", nil)
	w := httptest.NewRecorder()

	fakeJson := `{"test":"json"}`

	server.JwksHandler(w, r, fakeJson)
	result := w.Result()
	defer result.Body.Close()

	assert.Equal(t, http.StatusOK, result.StatusCode)

	body, err := ioutil.ReadAll(result.Body)
	assert.Nil(t, err, "should not have errored")
	assert.Equal(t, fakeJson, string(body))
}

func TestJWKSBadRequest(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/api/v1/jwks.json", nil)
	w := httptest.NewRecorder()

	server.JwksHandler(w, r, "")
	result := w.Result()

	assert.Equal(t, http.StatusBadRequest, result.StatusCode)
}
