package server_test

import (
	"io/ioutil"
	"jwks-server/internal/server"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
