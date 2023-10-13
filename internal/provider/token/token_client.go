package token

import (
	"GatewayService/internal/config"
	"fmt"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"time"
)

type JWTProvider struct {
	client       http.Client
	url          string
	retry        int
	timeoutRetry time.Duration
	logger       *zap.Logger
}

func NewJWTProvider(cfg config.JWTProvider, logger *zap.Logger) (*JWTProvider, error) {
	client := http.Client{
		Timeout: cfg.Timeout,
	}

	url := fmt.Sprintf("http://%s:%d", cfg.Host, cfg.Port)

	provider := &JWTProvider{
		client:       client,
		url:          url,
		retry:        cfg.Retry,
		timeoutRetry: cfg.TimeoutRetry,
		logger:       logger,
	}

	provider.logger.Info(provider.url)

	if err := provider.Ping(); err != nil {
		return nil, err
	}

	return provider, nil
}

func (provider *JWTProvider) Ping() error {
	_, err := Repeater(provider.retry, provider.timeoutRetry, func() (interface{}, error) {
		req, err := http.NewRequest("GET", provider.url+"/ping", nil)

		if err != nil {
			return nil, err
		}

		resp, err := provider.client.Do(req)

		if err != nil {
			return nil, err
		}

		if resp.StatusCode != 200 {
			return nil, fmt.Errorf("Cant dial %s with status code %d", resp.Request.URL, resp.StatusCode)
		}

		return nil, nil
	})

	return err
}

func (provider *JWTProvider) GetJWTToken(login string) (string, error) {
	token, err := Repeater(provider.retry, provider.timeoutRetry, func() (interface{}, error) {
		urlWithParams := fmt.Sprintf("%s/generate?login=%s", provider.url, login)
		req, err := http.NewRequest("GET", urlWithParams, nil)
		if err != nil {
			return nil, err
		}

		resp, err := provider.client.Do(req)

		if err != nil {
			return nil, err
		}

		if resp.StatusCode != 200 {
			if resp.StatusCode == 400 {
				return nil, fmt.Errorf("Token not found in header")
			}
			return nil, fmt.Errorf("Cant dial %s with status code %d", resp.Request.URL, resp.StatusCode)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		return body, nil
	})

	if err != nil {
		return "", err
	}

	b, ok := token.([]byte)
	if !ok {
		return "", fmt.Errorf("Generated invalid token %v", token)
	}
	tokenStr := string(b)

	return tokenStr, nil
}

func (provider *JWTProvider) ValidateToken(header string) error {

	_, err := Repeater(provider.retry, provider.timeoutRetry, func() (interface{}, error) {
		req, err := http.NewRequest("GET", provider.url+"/validate", nil)
		if err != nil {
			return nil, err
		}

		req.Header = http.Header{
			"Authorization": []string{"bearer " + header},
		}
		resp, err := provider.client.Do(req)

		if err != nil {
			return nil, err
		}

		if resp.StatusCode != 200 {
			if resp.StatusCode == 400 {
				return nil, fmt.Errorf("Token not found in header")
			} else if resp.StatusCode == 401 {
				return nil, fmt.Errorf("Invalid or expired token")
			}
			return nil, fmt.Errorf("Cant dial %s with status code %d", resp.Request.URL, resp.StatusCode)
		}

		return nil, nil
	})

	return err
}

type retryFunc func() (interface{}, error)

func Repeater(repeat int, timeoutEach time.Duration, exec retryFunc) (res interface{}, err error) {
	for i := 0; i < repeat; i++ {
		res, err = exec()
		if err == nil {
			return res, nil
		}
		time.Sleep(timeoutEach)
	}
	return nil, fmt.Errorf("Repeater max reapeat times")
}
