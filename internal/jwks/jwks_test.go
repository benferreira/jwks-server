package jwks_test

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"jwks-server/internal/config"
	"jwks-server/internal/jwks"
	"jwks-server/internal/rsa_helper"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewJWKSGeneration(t *testing.T) {
	keySet, err := jwks.NewJWKS(nil)

	assert.Nil(t, err, "should not have errored")
	assert.NotNil(t, keySet, "should have returned keyset")
	assert.Equal(t, 1, len(keySet.Keys))
	assert.Equal(t, "RS256", keySet.Keys[0].Alg)
	assert.Equal(t, "RSA", keySet.Keys[0].Kty)
	assert.NotEqual(t, "", keySet.Keys[0].N)
	assert.Equal(t, "AQAB", keySet.Keys[0].E)
}

func TestNewJWKSWithProvidedKey(t *testing.T) {
	key, _ := rsa_helper.GenerateRSAPublicKeyPem()
	pubKeys := []config.RSAPubKey{{Key: key}}

	keySet, err := jwks.NewJWKS(pubKeys)

	assert.Nil(t, err, "should not have errored")
	assert.NotNil(t, keySet, "should have returned keyset")
	assert.Equal(t, 1, len(keySet.Keys))
	assert.Equal(t, "RS256", keySet.Keys[0].Alg)
	assert.Equal(t, "RSA", keySet.Keys[0].Kty)
	assert.NotEqual(t, "", keySet.Keys[0].N)
	assert.Equal(t, "AQAB", keySet.Keys[0].E)
}

func TestNewJWK(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.Nil(t, err, "should have generated an RSA key")

	pubKeyString, err := marshalToPem(privateKey)
	assert.Nil(t, err, "should have marshalled key")

	jwk, err := jwks.NewJWK(config.RSAPubKey{Key: pubKeyString})
	assert.Nil(t, err, "should not have errored")
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
