package server

import (
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
)

func Serve(port int, jwksJson string) {
	http.HandleFunc("/health", HealthCheckHandler)

	http.HandleFunc("/api/v1/jwks.json", func(w http.ResponseWriter, r *http.Request) {
		JwksHandler(w, r, jwksJson)
	})

	log.Info().Msgf("serving localhost:%d", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func JwksHandler(w http.ResponseWriter, r *http.Request, jwksJson string) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, jwksJson)
}
