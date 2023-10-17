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

func NewServer(cfg *config.HTTPServerConfig, handler http.Handler, logger *zap.Logger) *Server {
	server := http.Server{
		Addr:              cfg.Host + ":" + cfg.Port,
		Handler:           handler,
		ReadTimeout:       cfg.ReadTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
	}
	return &Server{
		httpServer: &server,
		timeOutSec: cfg.TimeOutSec,
		logger:     logger,
	}
}

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
