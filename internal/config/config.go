package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
)

// HTTPServer holds config information for http server
type HTTPServer struct {
	MaxHeaderBytes    int
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	ReadHeaderTimeout time.Duration
	Port              string
	Host              string
}

// Server contains setting for setting up Server that contains HTTP and gRPC servers
type Server struct {
	HTTP       *HTTPServer
	TimeOutSec int
}

// Configurator used to get all configurations from configured conf files
type Configurator struct {
}

// NewConfiguration set paths to config.yml
func NewConfiguration() (*Configurator, error) {
	godotenv.Load("../../.env")
	switch os.Getenv("CURRENT_ENV") {
	case "local":

		viper.AddConfigPath("../../configs")
		viper.SetConfigName("config")
	default:

		viper.AddConfigPath("../../configs")
		viper.SetConfigName("config")
	}

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read conf file: %w", err)
	}

	c := &Configurator{}

	return c, nil
}

// SrvConfig returns configuration for server
func (cfg *Configurator) HTTPSrvConfig() *HTTPServer {
	log.Println("Reading configuration")

	return &HTTPServer{
		MaxHeaderBytes:    viper.GetInt("srv.maxHeaderBytes"),
		ReadTimeout:       viper.GetDuration("srv.readTimeout"),
		WriteTimeout:      viper.GetDuration("srv.writeTimeout"),
		ReadHeaderTimeout: viper.GetDuration("srv.readHeaderTimeout"),
		Port:              viper.GetString("srv.port"),
		Host:              viper.GetString("srv.host"),
	}
}

// ServerConfig returns config for creating Server struct
func (cfg *Configurator) ServerConfig() *Server {
	log.Println("Reading server config")
	httpSrv := cfg.HTTPSrvConfig()

	srv := Server{
		HTTP:       httpSrv,
		TimeOutSec: viper.GetInt("srv.timeOutSec"),
	}

	return &srv
}

type AppEnvs int

const (
	Release = iota
	Development
)

// GetEnvironment returns application development stage
func (cfg *Configurator) GetEnvironment(logger *zap.Logger) AppEnvs {
	logger.With(
		zap.String("place", "GetEnvironment"),
	).Info("Reading GetEnvironment")

	switch os.Getenv("APP_ENV") {
	case "release":
		logger.Info("Running in release")

		return Release
	default:
		logger.Info("Running in development")

		return Development
	}
}

type JWTProvider struct {
	Host         string
	Port         int
	Timeout      time.Duration
	Retry        int
	TimeoutRetry time.Duration
}

// JWTProviderConfig returns configuration for jwt generator
func (cfg *Configurator) JWTProviderConfig(logger *zap.Logger) *JWTProvider {
	logger.With(
		zap.String("place", "JWTProviderConfig"),
	).Info("Reading JWTProvider config from file")

	provider := &JWTProvider{
		Host:         viper.GetString("tokenGen.host"),
		Port:         viper.GetInt("tokenGen.port"),
		Timeout:      viper.GetDuration("tokenGen.timeout"),
		Retry:        viper.GetInt("tokenGen.retry"),
		TimeoutRetry: viper.GetDuration("tokenGen.timeoutRetry"),
	}
	return provider
}
