package jwks

import (
	"encoding/base64"
	"jwks-server/internal/config"
	"jwks-server/internal/rsa_helper"
)

type JWKS struct {
	Keys []JWK `json:"keys"`
}

// NewJWKS returns a JWKS with a JWK generated from the provided rsaPublicKey.
func NewJWKS(keys []config.RSAPubKey) (*JWKS, error) {
	if keys == nil {
		return generateJWKS()
	}

	jwks := JWKS{Keys: make([]JWK, 0)}

	for _, key := range keys {
		jwk, err := NewJWK(key)

		if err != nil {
			return nil, err
		}

		jwks.Keys = append(jwks.Keys, *jwk)
	}

	return &jwks, nil
}

func generateJWKS() (*JWKS, error) {
	generatedKey, err := rsa_helper.GenerateRSAPublicKeyPem()

	if err != nil {
		return nil, err
	}

	jwk, err := NewJWK(config.RSAPubKey{Key: generatedKey})

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
func NewJWK(rsaPublicKey config.RSAPubKey) (*JWK, error) {
	publicKey, err := rsa_helper.DecodeAndParsePKIXPublicKey(rsaPublicKey.Key)

	if err != nil {
		return nil, err
	}

	jwk := JWK{
		Alg: "RS256",
		Kty: "RSA",
		N:   base64.RawStdEncoding.EncodeToString(publicKey.N.Bytes()),
		E:   "AQAB",
	}

	if rsaPublicKey.Kid != "" {
		jwk.Kid = rsaPublicKey.Kid
	}

	return &jwk, nil
}
