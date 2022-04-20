package test_helper

import (
	"log"
	"net/http"
	"os"
	"time"
)

func UnsetTestEnvironment() {
	os.Unsetenv("DEBUG")
	os.Unsetenv("PORT")
	os.Unsetenv("PRETTY_LOGGING")
	os.Unsetenv("TEST_MODE")
	os.Unsetenv("TLS")
	os.Unsetenv("TLS_PRIVATE_KEY_PATH")
	os.Unsetenv("TLS_CERT_PATH")
	os.Unsetenv("RSA_PUB_KEY")
	os.Unsetenv("RSA_KEYS_FILE")
}

type TestHttpClient struct {
	*http.Client
	RetryLimit    int
	RetryInterval time.Duration
}

func NewTestHttpClient() *TestHttpClient {
	return &TestHttpClient{
		Client:        &http.Client{Timeout: time.Duration(1) * time.Second},
		RetryLimit:    3,
		RetryInterval: time.Duration(2) * time.Second,
	}
}

// GetWithRetry performs an HTTP GET with retry on failure up to TestHttpClient.RetryLimit
func (client *TestHttpClient) GetWithRetry(path string) (*http.Response, error) {
	req, err := http.NewRequest("GET", path, nil)

	if err != nil {
		return nil, err
	}

	return client.DoWithRetry(req)
}

// DoWithRetry performs an http.Client.Do() with retry on failure up to TestHttpClient.RetryLimit
func (client *TestHttpClient) DoWithRetry(req *http.Request) (*http.Response, error) {
	var res *http.Response
	var err error

	attempts := 1
	for res, err = client.Do(req); err != nil && attempts < client.RetryLimit; {
		log.Printf("error received, retrying in %v", client.RetryInterval)
		time.Sleep(client.RetryInterval)
		res, err = client.Do(req)
		attempts += 1
	}

	return res, err
}
