package jwks

import (
	"encoding/base64"
	"jwks-server/internal/rsa_helper"
)

type JWKS struct {
	Keys []JWK `json:"keys"`
}

// NewJWKS returns a JWKS with a JWK generated from the provided rsaPublicKey.
func NewJWKS(rsaPublicKey string) (*JWKS, error) {
	if rsaPublicKey != "" {
		jwk, err := NewJWK(rsaPublicKey)

		if err != nil {
			return nil, err
		}

		return &JWKS{Keys: []JWK{*jwk}}, nil
	}

	generatedKey, err := rsa_helper.GenerateRSAPublicKeyPem()

	if err != nil {
		return nil, err
	}

	jwk, err := NewJWK(generatedKey)

	if err != nil {
		return nil, err
	}

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

// NewJWK returns a JWK generated from the provided rsaPublicKey.
func NewJWK(rsaPublicKey string) (*JWK, error) {
	publicKey, err := rsa_helper.DecodeAndParsePKIXPublicKey(rsaPublicKey)

	if err != nil {
		return nil, err
	}

	jwk := JWK{
		Alg: "RS256",
		Kty: "RSA",
		N:   base64.RawURLEncoding.EncodeToString(publicKey.N.Bytes()),
		E:   "AQAB",
	}

	return &jwk, nil
}
