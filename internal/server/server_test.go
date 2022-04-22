package server_test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"math/big"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	test_helper "github.com/benferreira/jwks-server/_test_helper"
	"github.com/benferreira/jwks-server/internal/config"
	"github.com/benferreira/jwks-server/internal/server"

	"github.com/stretchr/testify/assert"
)

func TestServe(t *testing.T) {
	serv, err := server.NewServer(&config.ServerConfig{Port: 45566}, `{"test":"json"}`)
	assert.Nil(t, err)

	go func() {
		client := test_helper.NewTestHttpClient()
		resp, err := client.GetWithRetry("http://127.0.0.1:45566/health")
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
	defer result.Body.Close()

	assert.Equal(t, http.StatusOK, result.StatusCode)
	assert.Equal(t, "application/json", result.Header.Get("content-type"))

	body, err := io.ReadAll(result.Body)
	assert.Nil(t, err)
	assert.Equal(t, server.HealthCheckUpJson(), string(body))
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
	assert.Equal(t, "application/json", result.Header.Get("content-type"))
	assert.NotEmpty(t, result.Header.Get("cache-control"))

	body, err := io.ReadAll(result.Body)
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

func TestServeTLS(t *testing.T) {
	certPath, keyPath := generate_cert("../../_test_helper")
	defer delete_cert("../../_test_helper")
	conf := config.ServerConfig{
		Port:              45569,
		TLS:               true,
		TLSPrivateKeyPath: keyPath,
		TLSCertPath:       certPath,
	}

	serv, err := server.NewServer(&conf, `{"test":"json"}`)
	assert.Nil(t, err)

	go func() {
		client := test_helper.NewTestHttpClient()
		client.Transport = &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}

		resp, err := client.GetWithRetry("https://127.0.0.1:45569/health")
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		serv.Shutdown(context.Background())
	}()

	err = serv.Start()
	assert.Nil(t, err)
}

// generates an RSA private key and x509 certificate
// code taken from https://go.dev/src/crypto/tls/generate_cert.go
func generate_cert(path string) (string, string) {
	keyUsage := x509.KeyUsageDigitalSignature
	keyUsage |= x509.KeyUsageKeyEncipherment

	notBefore := time.Now()
	notAfter := notBefore.Add(time.Minute * 1)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		log.Fatalf("Failed to generate serial number: %v", err)
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Acme Co"},
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,

		KeyUsage:              keyUsage,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,

		IsCA: true,
	}

	template.KeyUsage |= x509.KeyUsageCertSign

	hosts := []string{"127.0.0.1", "::1", "localhost"}
	for _, h := range hosts {
		if ip := net.ParseIP(h); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, h)
		}
	}

	priv, err := rsa.GenerateKey(rand.Reader, 2048)

	if err != nil {
		log.Fatalf("Failed to generate rsa key: %v", err)
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		log.Fatalf("Failed to create certificate: %v", err)
	}

	certOut, err := os.Create(fmt.Sprintf("%s/%s", path, "cert.pem"))
	if err != nil {
		log.Fatalf("Failed to open cert.pem for writing: %v", err)
	}
	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		log.Fatalf("Failed to write data to cert.pem: %v", err)
	}
	if err := certOut.Close(); err != nil {
		log.Fatalf("Error closing cert.pem: %v", err)
	}
	log.Print("wrote cert.pem\n")

	keyOut, err := os.OpenFile(fmt.Sprintf("%s/%s", path, "key.pem"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Failed to open key.pem for writing: %v", err)
	}
	privBytes, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		log.Fatalf("Unable to marshal private key: %v", err)
	}
	if err := pem.Encode(keyOut, &pem.Block{Type: "PRIVATE KEY", Bytes: privBytes}); err != nil {
		log.Fatalf("Failed to write data to key.pem: %v", err)
	}
	if err := keyOut.Close(); err != nil {
		log.Fatalf("Error closing key.pem: %v", err)
	}
	log.Print("wrote key.pem\n")

	return certOut.Name(), keyOut.Name()
}

func delete_cert(path string) {
	os.Remove(fmt.Sprintf("%s/%s", path, "cert.pem"))
	os.Remove(fmt.Sprintf("%s/%s", path, "key.pem"))
}
