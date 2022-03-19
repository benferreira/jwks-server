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
		jwksHandler(w, r, jwksJson)
	})

	log.Info().Msgf("serving localhost:%d", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func jwksHandler(w http.ResponseWriter, r *http.Request, jwksJson string) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, jwksJson)
}
