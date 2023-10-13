package main

import (
	"GatewayService/internal/config"
	"GatewayService/internal/handler"
	"GatewayService/internal/middleware"
	"GatewayService/internal/provider/token"
	"GatewayService/internal/repository"
	"GatewayService/internal/server"
	"GatewayService/internal/service"
	"context"
	"fmt"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg, err := config.NewConfiguration()

	if err != nil {
		log.Fatalf("Failed to initialize config: %v", err)
	}

	logger, err := initLogger()

	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	//add information about environment into logger
	isRelease := cfg.GetEnvironment(logger) != config.Development
	envField := zap.String("APP_ENV", "production")
	if !isRelease {
		envField = zap.String("APP_ENV", "development")
	}

	logger.Info("Got application environment", envField)

	providerCfg := cfg.JWTProviderConfig(logger)

	authProvider, err := token.NewJWTProvider(*providerCfg, logger)
	if err != nil {
		logger.With(
			zap.String("place", "system(main)"),
			zap.Error(err),
		).Panic("Failed to connect to token generator")
	}

	userRepository := repository.NewMockUserRepository()

	authService := service.NewAuthService(authProvider, logger, userRepository)

	authHandler := handler.NewAuthHandler(authService, logger)
	shopsHandler := handler.NewShopsHandler(logger)

	authMiddleware := middleware.NewMiddleware(authProvider)

	router := handler.NewRouter(authHandler, shopsHandler, authMiddleware)

	srvCfg := cfg.ServerConfig()
	srv := server.NewServer(srvCfg, router)
	fmt.Println(srv)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		if err := srv.Run(ctx); err != nil {
			logger.With(
				zap.String("place", "system(main)"),
				zap.Error(err),
			).Error("Server failed during run")
		}
	}()

	shutDown := make(chan os.Signal, 1)

	signal.Notify(shutDown, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	s := <-shutDown

	logger.With(
		zap.String("signal", s.String()),
	).Info("System received signal for shutdown, shutting down server")

	cancel()
}

func initLogger() (*zap.Logger, error) {
	logger, err := zap.NewProduction()
	if os.Getenv("APP_ENV") == "development" {
		logger, err = zap.NewDevelopment()
	}
	return logger, err
}
