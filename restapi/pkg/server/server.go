// Package server contains custom http server for our API
package server // import "github.com/la4ezar/restapi/pkg/server

import (
	"context"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/la4ezar/restapi/internal/crypto"
	"github.com/la4ezar/restapi/pkg/controller"
	"github.com/la4ezar/restapi/pkg/log"

	"github.com/gorilla/mux"
)

type Cryptocurrency crypto.Cryptocurrency
type CryptoAuthors []crypto.Author

// Server is wrapped http.Server with shutdown timeout
type Server struct {
	*http.Server

	shutdownTimeout time.Duration
}

// New returns new Server instance with given configurations and router
func New(cfg *Config, ctr controller.Controller) *Server {
	r := mux.NewRouter().StrictSlash(true)

	r.Use(log.RequestLogger())

	for _, route := range *ctr.Routes() {
		r.Name(route.Name).
			Methods(route.Method).
			Path(route.Path).
			Handler(route.Handler)
	}

	s := &Server{
		Server: &http.Server{
			Addr:         ":" + strconv.Itoa(cfg.Port),
			Handler:      r,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
			IdleTimeout:  cfg.IdleTimeout,
		},
		shutdownTimeout: cfg.ShutdownTimeout,
	}

	return s
}

// Start starts the Server
func (s *Server) Start(parentCtx context.Context, wg *sync.WaitGroup) {
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	// Checks if the ctx or parentCtx is canceled
	go func() {
		defer wg.Done()
		<-ctx.Done()
		s.stop()
	}()

	log.C(ctx).Infof("Starting and listening on %s\n", s.Addr)

	log.C(ctx).Fatal(s.ListenAndServe())
}

// stop stops the Server
func (s *Server) stop() {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	s.SetKeepAlivesEnabled(false)

	if err := s.Shutdown(ctx); err != nil {
		log.C(ctx).Fatalf("Couldn't gracefully shutdown the server: %v\n", err)
	}

	log.C(ctx).Infof("Server gracefully shutdown.")
}
