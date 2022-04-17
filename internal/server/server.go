package server

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/benferreira/jwks-server/internal/config"
	"github.com/rs/zerolog/log"
)

type Server struct {
	*http.Server
	*config.ServerConfig
}

// NewServer constructs a new http.Server and registers the handlers
func NewServer(conf *config.ServerConfig, jwksJson string) (*Server, error) {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", HealthCheckHandler)

	mux.HandleFunc("/api/v1/jwks.json", func(w http.ResponseWriter, r *http.Request) {
		JwksHandler(w, r, jwksJson)
	})

	server := &Server{
		&http.Server{
			Addr:    fmt.Sprintf(":%d", conf.Port),
			Handler: mux,
		},
		conf,
	}

	if conf.TLS {
		cert, err := tls.LoadX509KeyPair(conf.TLSCertPath, conf.TLSPrivateKeyPath)
		if err != nil {
			return nil, err
		}

		server.TLSConfig = &tls.Config{Certificates: []tls.Certificate{cert}}
	}

	return server, nil
}

func (serv *Server) Start() error {
	log.Info().Msgf("serving localhost%s", serv.Addr)

	var err error

	if serv.TLS {
		log.Info().Msg("TLS enabled")
		err = serv.ListenAndServeTLS("", "")
	} else {
		err = serv.ListenAndServe()
	}

	if err != http.ErrServerClosed {
		return err
	} else {
		return nil
	}
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
