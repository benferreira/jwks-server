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
	RsaPubKeys []RSAPubKey
	TestMode   bool
}

type RSAPubKey struct {
	Key string
	Kid string
}

func NewConfig() (*Config, error) {
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
		conf.RsaPubKeys = make([]RSAPubKey, 0)
		conf.RsaPubKeys = append(conf.RsaPubKeys, RSAPubKey{Key: rsaKey})
	} else {
		return nil, fmt.Errorf("RSA_PUB_KEY env var must be provided")
	}

	return &conf, nil
}
