package jwks_test

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"testing"

	"github.com/benferreira/jwks-server/internal/config"
	"github.com/benferreira/jwks-server/internal/jwks"
	"github.com/benferreira/jwks-server/internal/rsa_helper"

	"github.com/MicahParks/keyfunc"
	"github.com/golang-jwt/jwt/v4"
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
	pubKeys, _ := config.NewRSAPubKeys(key)

	keySet, err := jwks.NewJWKS(pubKeys)

	assert.Nil(t, err, "should not have errored")
	assert.NotNil(t, keySet, "should have returned keyset")
	assert.Equal(t, 1, len(keySet.Keys))
	assert.Equal(t, "RS256", keySet.Keys[0].Alg)
	assert.Equal(t, "RSA", keySet.Keys[0].Kty)
	assert.NotEqual(t, "", keySet.Keys[0].N)
	assert.Equal(t, "AQAB", keySet.Keys[0].E)
}

func TestNewJWKSWIthMultipleKeys(t *testing.T) {
	key1, _ := rsa_helper.GenerateRSAPublicKeyPem()
	key2, _ := rsa_helper.GenerateRSAPublicKeyPem()

	pubKeys, _ := config.NewRSAPubKeys(key1)
	pubKeys.RSAPubKeys[0].Kid = "key1"
	pubKeys.RSAPubKeys = append(pubKeys.RSAPubKeys, config.RSAPubKey{Key: key2, Kid: "key2"})

	keySet, err := jwks.NewJWKS(pubKeys)
	assert.Nil(t, err, "should have generated a keyset")
	assert.NotNil(t, keySet, "should have generated a keyset")
	assert.Equal(t, "key1", keySet.Keys[0].Kid)
	assert.Equal(t, "key2", keySet.Keys[1].Kid)

	for _, key := range keySet.Keys {
		assert.NotNil(t, key.N)
	}
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

// TestJWKSisValid generates a new JWT we sign ourselves with a new RSA key. We then verify that our JWKS generated
// from the RSA key can be used to successfully validate the JWT.
func TestJWKSisValid(t *testing.T) {
	//Generate JWT and sign it with new private key
	token := jwt.New(jwt.SigningMethodRS256)
	token.Header["kid"] = "kid1"
	privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	tokenString, err := token.SignedString(privateKey)
	assert.Nil(t, err)
	assert.NotNil(t, tokenString)

	// populate a JWKS with public key
	jwksJson := generateJWKSJson(privateKey, "kid1")

	// Take our jwks and provide it to the JWT parser
	jwksKf, err := keyfunc.NewJSON(jwksJson)
	assert.Nil(t, err)

	parsedToken, err := jwt.Parse(tokenString, jwksKf.Keyfunc)
	assert.Nil(t, err, "token failed to parse")
	assert.True(t, parsedToken.Valid, "token should be valid")
}

func generateJWKSJson(privateKey *rsa.PrivateKey, kid string) []byte {
	publicKeyString, _ := marshalToPem(privateKey)
	pubKey := config.RSAPubKey{Key: publicKeyString, Kid: kid}
	jwks, _ := jwks.NewJWKS(&config.RSAPubKeys{RSAPubKeys: []config.RSAPubKey{pubKey}})
	jwksJson, _ := json.Marshal(jwks)

	return jwksJson
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
