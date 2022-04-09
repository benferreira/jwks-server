package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Debug      bool
	Port       int
	PrettyLog  bool
	PublicKeys *RSAPubKeys
	TestMode   bool
}

func NewConfigFromEnv() (*Config, error) {
	conf := Config{
		Debug:    false,
		Port:     8000,
		TestMode: false,
	}

	if _, ok := os.LookupEnv("DEBUG"); ok {
		conf.Debug = true
	}

	if port, ok := os.LookupEnv("PORT"); ok {
		portInt, err := strconv.Atoi(port)

		if err != nil {
			return nil, fmt.Errorf("invalid value for PORT: error: %v", err)
		}

		conf.Port = portInt
	}

	if _, ok := os.LookupEnv("PRETTY_LOGGING"); ok {
		conf.PrettyLog = true
	}

	if _, ok := os.LookupEnv("TEST_MODE"); ok {
		conf.TestMode = true
		return &conf, nil
	}

	if rsaKey, ok := os.LookupEnv("RSA_PUB_KEY"); ok {
		if keys, err := NewRSAPubKeys(rsaKey); err == nil {
			conf.PublicKeys = keys
		} else {
			return nil, err
		}
	} else if path, ok := os.LookupEnv("RSA_KEYS_FILE"); ok {
		if keys, err := keysFromFile(path); err == nil {
			conf.PublicKeys = keys
		} else {
			return nil, err
		}
	}

	if err := conf.validate(); err != nil {
		return nil, err
	}

	return &conf, nil
}

func keysFromFile(path string) (*RSAPubKeys, error) {
	if file, err := os.Open(path); err == nil {
		return NewRSAPubKeysFromFile(file)
	} else {
		return nil, err
	}
}

func (c Config) validate() error {

	if c.Port == 0 {
		return fmt.Errorf("invalid configuration, port must be set")
	}

	//If test mode is enabled, RsaPubKey is not required
	if c.TestMode {
		return nil
	}

	if c.PublicKeys == nil {
		return fmt.Errorf("invalid configuration, missing public keys")
	}

	return nil
}
