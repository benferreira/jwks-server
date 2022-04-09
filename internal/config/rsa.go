package config

import (
	"fmt"
	"io"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type RSAPubKeys struct {
	RSAPubKeys []RSAPubKey `yaml:"rsaPubKeys"`
}

type RSAPubKey struct {
	Key string `yaml:"key"`
	Kid string `yaml:"kid"`
}

func NewRSAPubKeys(pubKey string) (*RSAPubKeys, error) {
	key := RSAPubKey{Key: pubKey}

	if err := key.validate(); err != nil {
		return nil, err
	}

	return &RSAPubKeys{[]RSAPubKey{key}}, nil
}

func NewRSAPubKeysFromFile(reader io.Reader) (*RSAPubKeys, error) {
	keysBytes, err := ioutil.ReadAll(reader)

	if err != nil {
		return nil, err
	}

	keys := RSAPubKeys{}
	if err = yaml.Unmarshal(keysBytes, &keys); err != nil {
		return nil, err
	}

	for _, key := range keys.RSAPubKeys {
		if err = key.validate(); err != nil {
			return nil, err
		}
	}

	return &keys, nil
}

func (r RSAPubKey) validate() error {
	if r.Key == "" {
		return fmt.Errorf("invalid configuration, public key value is missing")
	}

	return nil
}
