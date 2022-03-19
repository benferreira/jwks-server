package jwks

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
)

type JWKS struct {
	Keys []JWK `json:"keys"`
}

func NewJWKS(rsaPublicKey string) (*JWKS, error) {
	if rsaPublicKey != "" {
		jwk, err := NewJWK(rsaPublicKey)

		if err != nil {
			return nil, err
		}

		return &JWKS{Keys: []JWK{*jwk}}, nil
	}

	generatedKey, err := generateRSAPublicKeyPem()

	if err != nil {
		return nil, err
	}

	jwk, err := NewJWK(generatedKey)

	return &JWKS{Keys: []JWK{*jwk}}, nil
}

type JWK struct {
	Alg string   `json:"alg"`
	Kty string   `json:"kty"`
	Use string   `json:"use,omitempty"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	Kid string   `json:"kid,omitempty"`
	X5t string   `json:"x5t,omitempty"`
	X5c []string `json:"x5c,omitempty"`
}

func NewJWK(rsaPub string) (*JWK, error) {
	block, _ := pem.Decode([]byte(rsaPub))

	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, fmt.Errorf("failed to decode PEM block containing public key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse PKIX public key")
	}

	publicKey := pub.(*rsa.PublicKey)

	jwk := JWK{
		Alg: "RS256",
		Kty: "RSA",
		N:   base64.RawStdEncoding.EncodeToString(publicKey.N.Bytes()),
		E:   "AQAB",
	}

	return &jwk, nil
}

func generateRSAPublicKeyPem() (string, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)

	if err != nil {
		return "", err
	}

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
