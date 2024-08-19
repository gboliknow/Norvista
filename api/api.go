package api

import (
	"net/http"
	"os"

	"github.com/rs/zerolog"

	"github.com/gorilla/mux"
)

type APIServer struct {
	addr   string
	store  Store
	logger zerolog.Logger
}

func NewAPIServer(addr string, store Store) *APIServer {
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()
	return &APIServer{addr: addr, store: store, logger: logger}
}

func (s *APIServer) Serve() {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	//registering the routes

	s.logger.Info().Str("addr", s.addr).Msg("Starting API server")
	if err := http.ListenAndServe(s.addr, subrouter); err != nil {
		s.logger.Fatal().Err(err).Msg("Server stopped")
	}
}
