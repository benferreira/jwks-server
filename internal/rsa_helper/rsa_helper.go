package rsa_helper

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"github.com/rs/zerolog/log"
)

// GenerateRSAPublicKeyPem generates a new random RSA public key marshalled to a PKIX pem.
func GenerateRSAPublicKeyPem() (string, error) {
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

	log.Debug().Msgf("generated public key: %s", buf.String())

	return buf.String(), nil
}

// DecodeAndParsePKIXPublicKey takes an RSA public key, decodes it and parses it as PKIX.
func DecodeAndParsePKIXPublicKey(rsaPublicKey string) (*rsa.PublicKey, error) {
	block, rest := pem.Decode([]byte(rsaPublicKey))

	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, fmt.Errorf("failed to decode PEM block containing public key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse PKIX public key")
	}

	log.Debug().Msgf("parsed %T with remaining data: %q", pub, rest)

	return pub.(*rsa.PublicKey), nil
}
