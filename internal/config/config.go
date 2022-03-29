package config

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Debug      bool        `yaml:"debug"`
	Port       int         `yaml:"port"`
	PrettyLog  bool        `yaml:"prettyLogging"`
	RsaPubKeys []RSAPubKey `yaml:"rsaPubKeys"`
	TestMode   bool        `yaml:"testMode"`
}

type RSAPubKey struct {
	Key string `yaml:"key"`
	Kid string `yaml:"kid"`
}

func (c Config) ok() error {
	if c.TestMode {
		return nil
	}

	if c.RsaPubKeys == nil || len(c.RsaPubKeys) == 0 {
		return fmt.Errorf("invalid configuration, missing public key(s)")
	}

	return nil
}

func (r RSAPubKey) ok() error {
	if r.Key == "" {
		return fmt.Errorf("invalid configuration, public key value is missing")
	}

	return nil
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

func NewConfigFromFile(reader io.Reader) (*Config, error) {
	configBytes, err := ioutil.ReadAll(reader)

	if err != nil {
		return nil, err
	}

	config := Config{}
	if err = yaml.Unmarshal(configBytes, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
