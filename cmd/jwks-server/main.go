package main

import (
	"encoding/json"
	"flag"
	"jwks-server/internal/jwks"
	"jwks-server/internal/server"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var debug = flag.Bool("d", false, "enable debug logging")
var port = flag.Int("p", 8000, "http listening port")
var test = flag.Bool("t", false, "generate an RSA key and serve a JWKS")

func main() {
	parseFlags()

	if !*test && os.Getenv("RSA_PUB_KEY") == "" {
		log.Fatal().Msg("RSA_PUB_KEY env var must be provided")
	}

	constructedJWKS, err := buildJWKS()

	if err != nil {
		log.Fatal().Err(err).Msg("failed to construct JWKS")
	}

	jwksJson, err := json.Marshal(constructedJWKS)

	if err != nil {
		log.Fatal().Err(err).Msg("failed to marshall jwks to JSON")
	}

	server.Serve(*port, string(jwksJson))
}

func buildJWKS() (*jwks.JWKS, error) {
	if *test {
		return jwks.NewJWKS("")
	}

	return jwks.NewJWKS(os.Getenv("RSA_PUB_KEY"))
}

func parseFlags() {
	flag.Parse()

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	//Set logging level
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
}
