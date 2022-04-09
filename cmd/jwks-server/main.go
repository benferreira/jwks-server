package main

import (
	"encoding/json"
	"jwks-server/internal/config"
	"jwks-server/internal/jwks"
	"jwks-server/internal/server"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	conf, err := config.NewConfigFromEnv()

	if err != nil {
		log.Fatal().Err(err).Msg("invalid configuration")
	}

	configureLogging(conf.Debug, conf.PrettyLog)

	constructedJWKS, err := buildJWKS(conf)

	if err != nil {
		log.Fatal().Err(err).Msg("failed to construct JWKS")
	}

	jwksJson, err := json.Marshal(constructedJWKS)

	if err != nil {
		log.Fatal().Err(err).Msg("failed to marshall jwks to JSON")
	}

	serv := server.NewServer(conf.Port, string(jwksJson))

	if err := serv.Start(); err != nil {
		log.Fatal().Err(err).Msg("server error")
	}
}

func buildJWKS(conf *config.Config) (*jwks.JWKS, error) {
	if conf.TestMode {
		log.Info().Msg("test mode enabled, generating RSA public key")
		return jwks.NewJWKS(nil)
	}

	return jwks.NewJWKS(conf.PublicKeys)
}

func configureLogging(debug bool, pretty bool) {
	//Set logging level
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	if pretty {
		log.Logger = log.Output(zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: "2006-01-02T15:04:05.999Z07:00",
		})
	}
}
