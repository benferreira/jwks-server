package server

import (
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
)

func Serve(port int, jwksJson string) {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/api/v1/jwks.json", func(w http.ResponseWriter, r *http.Request) {
		jwksRoute(w, r, jwksJson)
	})

	log.Info().Msgf("serving 127.0.0.1:%d", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func jwksRoute(w http.ResponseWriter, r *http.Request, jwksJson string) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(jwksJson))
}
