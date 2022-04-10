package server

import (
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
)

type Server struct {
	*http.Server
}

// NewServer constructs a new http.Server and registers the handlers
func NewServer(port int, jwksJson string) *Server {
	mux := http.NewServeMux()
	serv := Server{
		&http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: mux,
		},
	}

	mux.HandleFunc("/health", HealthCheckHandler)

	mux.HandleFunc("/api/v1/jwks.json", func(w http.ResponseWriter, r *http.Request) {
		JwksHandler(w, r, jwksJson)
	})
	return &serv
}

func (serv *Server) Start() error {
	log.Info().Msgf("serving localhost%s", serv.Addr)

	if err := serv.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}

	return nil
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, HealthCheckUpJson())
}

func JwksHandler(w http.ResponseWriter, r *http.Request, jwksJson string) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=15, stale-while-revalidate=15, stale-if-error=86400")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, jwksJson)
}
