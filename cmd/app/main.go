package main

import (
	"GatewayService/internal/config"
	"GatewayService/internal/handler"
	"GatewayService/internal/handler/error/mapper"
	"GatewayService/internal/middleware"
	"GatewayService/internal/provider/token"
	"GatewayService/internal/repository"
	"GatewayService/internal/server"
	"GatewayService/internal/service"
	"context"
	"github.com/streadway/amqp"
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
			zap.String("place", "main"),
			zap.Error(err),
		).Panic("Failed to connect to JWT provider")
	}

	rabbitChannel, err := initRabbitMQConnection(cfg, logger)

	if err != nil {
		logger.With(
			zap.String("place", "main"),
			zap.Error(err),
		).Panic("Failed to establish rabbitMQ connection")
	}

	userRepository := repository.NewMockUserRepository()

	authService := service.NewAuthService(authProvider, logger, userRepository)

	errorMapper := mapper.NewAuthErrorMapper()

	authHandler := handler.NewAuthHandler(authService, logger, errorMapper)

	storesHandler := handler.NewStoresHandler(rabbitChannel, logger)

	authMiddleware := middleware.NewMiddleware(authProvider)

	router := handler.NewRouter(authHandler, storesHandler, authMiddleware)

	srvCfg := cfg.ServerConfig()

	srv := server.NewServer(srvCfg, router, logger)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		if err := srv.Run(ctx); err != nil {
			logger.With(
				zap.String("place", "main"),
				zap.Error(err),
			).Error("Server failed during run")
		}
	}()

	//graceful shutdown using buffered channel
	shutDown := make(chan os.Signal, 1)

	signal.Notify(shutDown, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	s := <-shutDown

	logger.With(
		zap.String("signal", s.String()),
	).Info("Shutting down server")

	cancel()
}

func initLogger() (*zap.Logger, error) {
	logger, err := zap.NewProduction()
	if os.Getenv("APP_ENV") == "development" {
		logger, err = zap.NewDevelopment()
	}
	return logger, err
}

func initRabbitMQConnection(cfg *config.Configurator, logger *zap.Logger) (*amqp.Channel, error) {
	mqConfig := cfg.GetRabbitMQConfig()

	conn, err := amqp.Dial(cfg.GetAMQPConnectionURL(mqConfig))
	if err != nil {
		logger.With(
			zap.String("place", "initRabbitMQConnection"),
			zap.Error(err),
		).Error("Failed to establish RabbitMQ connection")
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		logger.With(
			zap.String("place", "initRabbitMQConnection"),
			zap.Error(err),
		).Error("Failed to open RabbitMQ channel")
		return nil, err
	}

	return channel, nil
}
