package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"os"
	"time"

	"github.com/spf13/viper"
)

type RabbitMQConfig struct {
	Host     string
	Port     string
	Username string
	Password string
}

// HTTPServerConfig holds config information for http server
type HTTPServerConfig struct {
	MaxHeaderBytes    int
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	ReadHeaderTimeout time.Duration
	TimeOutSec        int
	Port              string
	Host              string
}

// Server contains setting for setting up Server that contains HTTP and gRPC servers
type Server struct {
	HTTP *HTTPServerConfig
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
func (cfg *Configurator) getHTTPSrvConfig() *HTTPServerConfig {

	return &HTTPServerConfig{
		MaxHeaderBytes:    viper.GetInt("srv.maxHeaderBytes"),
		ReadTimeout:       viper.GetDuration("srv.readTimeout"),
		WriteTimeout:      viper.GetDuration("srv.writeTimeout"),
		ReadHeaderTimeout: viper.GetDuration("srv.readHeaderTimeout"),
		TimeOutSec:        viper.GetInt("srv.timeOutSec"),
		Port:              viper.GetString("srv.port"),
		Host:              viper.GetString("srv.host"),
	}
}

// ServerConfig returns config for creating Server struct
func (cfg *Configurator) ServerConfig() *Server {
	httpSrv := cfg.getHTTPSrvConfig()

	srv := Server{
		HTTP: httpSrv,
	}

	return &srv
}

func (cfg *Configurator) GetRabbitMQConfig() *RabbitMQConfig {
	return &RabbitMQConfig{
		Password: viper.GetString("rabbit.password"),
		Username: viper.GetString("rabbit.username"),
		Port:     viper.GetString("rabbit.port"),
		Host:     viper.GetString("rabbit.host"),
	}
}

func (cfg *Configurator) GetAMQPConnectionURL(rabbitCfg *RabbitMQConfig) string {
	return fmt.Sprintf("amqp://%s:%s@%s:%s/", rabbitCfg.Username, rabbitCfg.Password, rabbitCfg.Host, rabbitCfg.Port)
}

type AppEnvironment string

const (
	Release             AppEnvironment = "release"
	Development         AppEnvironment = "development"
	DefaultEnv          AppEnvironment = Development
	EnvironmentVariable                = "APP_ENV"
)

func (cfg *Configurator) GetEnvironment(logger *zap.Logger) AppEnvironment {
	logger.With(
		zap.String("place", "GetEnvironment"),
	).Info("Reading GetEnvironment")

	env := os.Getenv(EnvironmentVariable)
	if env == "" {
		env = string(DefaultEnv)
	}

	logger.Info("Running in " + env)
	return AppEnvironment(env)
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
