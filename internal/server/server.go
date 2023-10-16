package server

import (
	"GatewayService/internal/config"
	"context"
	"go.uber.org/zap"
	"net/http"
	"time"

	"golang.org/x/sync/errgroup"
)

type Server struct {
	httpServer *http.Server
	timeOutSec int
	logger     *zap.Logger
}

// function called to create HTTP server and configure it
func NewServer(cfg *config.Server, handler http.Handler, logger *zap.Logger) *Server {
	server := http.Server{
		Addr:              cfg.HTTP.Host + ":" + cfg.HTTP.Port,
		Handler:           handler,
		MaxHeaderBytes:    cfg.HTTP.MaxHeaderBytes,
		ReadTimeout:       cfg.HTTP.ReadTimeout,
		WriteTimeout:      cfg.HTTP.WriteTimeout,
		ReadHeaderTimeout: cfg.HTTP.ReadHeaderTimeout,
	}
	return &Server{
		httpServer: &server,
		timeOutSec: cfg.HTTP.TimeOutSec,
		logger:     logger,
	}
}

// Run starts configured server
// receives context for cancellation
// returns either error occurred during run, or during shutdown, or nothing
func (s *Server) Run(ctx context.Context) error {
	g, _ := errgroup.WithContext(ctx)
	g.Go(func() error {
		s.logger.Info("Server is running")
		return s.httpServer.ListenAndServe()
	})

	g.Go(func() error {
		<-ctx.Done()
		return s.Shutdown()
	})

	return g.Wait()
}

func (s *Server) Shutdown() error {
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Duration(s.timeOutSec)*time.Second)
	defer cancel()

	return s.httpServer.Shutdown(ctxTimeout)
}
