package jwks_test

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"jwks-server/internal/jwks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewJWK(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.Nil(t, err, "should have generated an RSA key")

	pubKeyString, err := marshalToPem(privateKey)
	assert.Nil(t, err, "should have marshalled key")

	jwk, err := jwks.NewJWK(pubKeyString)
	assert.Nil(t, err, "should have generated jwk")
	assert.NotNil(t, jwk, "should have generated jwk")
}

func marshalToPem(privateKey *rsa.PrivateKey) (string, error) {
	pubKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return "", err
	}

	pubKeyBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubKeyBytes,
	}

	buf := new(bytes.Buffer)
	pem.Encode(buf, pubKeyBlock)

	return buf.String(), nil
}
