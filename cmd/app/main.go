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

	if isRelease := cfg.GetEnvironment(logger) == config.Release; isRelease {
		logger.Info("Got application environment. Running in Release")
	} else {
		logger.Info("Got application environment. Running in Development")
	}

	providerCfg := cfg.JWTProviderConfig(logger)

	authProvider, err := token.NewJWTProvider(*providerCfg, logger)
	if err != nil {
		logger.With(
			zap.String("place", "main"),
			zap.Error(err),
		).Panic("Failed to connect to JWT provider")
	}

	rabbitConnection, err := initRabbitMQConnection(cfg)

	if err != nil {
		rabbitConnection.Close()
		logger.With(
			zap.String("place", "main"),
			zap.Error(err),
		).Panic("Failed to establish RabbitMQ rabbitConnection")
	}

	channel, err := initRabbitChannel(rabbitConnection)

	if err != nil {
		channel.Close()
		logger.With(
			zap.String("place", "main"),
			zap.Error(err),
		).Panic("Failed to init RabbitMQ channel")
	}

	queueName, err := declareRabbitQueue(channel)

	if err != nil {
		logger.With(
			zap.String("place", "main"),
			zap.Error(err),
		).Panic("Failed to init RabbitMQ queue")
	}

	userRepository := repository.NewMockUserRepository()

	authService := service.NewAuthService(authProvider, logger, userRepository)

	errorMapper := mapper.NewAuthErrorMapper()

	authHandler := handler.NewAuthHandler(authService, logger, errorMapper)

	storesHandler := handler.NewStoresHandler(channel, rabbitConnection, queueName, logger)

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

func declareRabbitQueue(channel *amqp.Channel) (string, error) {
	queue, err := channel.QueueDeclare(
		"CreateQueue", // name
		false,         // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	return queue.Name, err
}

func initRabbitChannel(connection *amqp.Connection) (*amqp.Channel, error) {
	channel, err := connection.Channel()

	return channel, err
}

func initRabbitMQConnection(cfg *config.Configurator) (*amqp.Connection, error) {
	mqConfig := cfg.GetRabbitMQConfig()

	conn, err := amqp.Dial(cfg.GetAMQPConnectionURL(mqConfig))

	return conn, err
}
