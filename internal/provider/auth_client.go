package provider

import (
	"GatewayService/internal/config"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
	"time"
)

type AuthProvider struct {
	client       http.Client
	url          string
	retry        int
	retryTimeout time.Duration
	logger       *zap.Logger
}

func NewAuthProvider(cfg config.AuthProviderConfig, logger *zap.Logger) (*AuthProvider, error) {
	client := http.Client{
		Timeout: cfg.Timeout,
	}

	url := fmt.Sprintf("http://%s:%d", cfg.Host, cfg.Port)

	provider := &AuthProvider{
		client:       client,
		url:          url,
		retry:        cfg.Retry,
		retryTimeout: cfg.TimeoutRetry,
		logger:       logger,
	}

	if err := provider.Ping(); err != nil {
		return nil, err
	}

	return provider, nil
}

func (p *AuthProvider) GetJWTToken(login string) (string, error) {
	urlWithParams := fmt.Sprintf("%s/generate?login=%s", p.url, login)
	req, err := http.NewRequest("GET", urlWithParams, nil)
	if err != nil {
		return "", err
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusBadRequest {
			return "", fmt.Errorf("token not found in header")
		}
		return "", fmt.Errorf("failed to make request to %s with status code %d", resp.Request.URL, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	tokenStr := string(body)
	return tokenStr, nil
}

func (p *AuthProvider) ValidateToken(header string) error {
	req, err := http.NewRequest("GET", p.url+"/validate", nil)
	if err != nil {
		return err
	}

	req.Header = http.Header{
		"Authorization": []string{"bearer " + header},
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusOK {
		return nil
	} else if resp.StatusCode == http.StatusBadRequest {
		return fmt.Errorf("token not found in header")
	} else if resp.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("invalid or expired token")
	}

	return fmt.Errorf("failed to validate token with status code %d", resp.StatusCode)
}

func RetryConnection(repeat int, timeoutEach time.Duration, exec func() (interface{}, error)) (res interface{}, err error) {
	for i := 0; i < repeat; i++ {
		res, err = exec()
		if err == nil {
			return res, nil
		}
		time.Sleep(timeoutEach)
	}

	return nil, fmt.Errorf("retry connection failed")
}

func (p *AuthProvider) Ping() error {
	_, err := RetryConnection(p.retry, p.retryTimeout, func() (interface{}, error) {
		req, err := http.NewRequest("GET", p.url+"/ping", nil)
		if err != nil {
			return nil, err
		}

		resp, err := p.client.Do(req)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("cant dial %s with status code %d", resp.Request.URL, resp.StatusCode)
		}

		return nil, nil
	})

	return err
}
