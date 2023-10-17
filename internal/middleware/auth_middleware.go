package middleware

import (
	"GatewayService/internal/handler/response"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"strings"
)

const (
	Header = "Authorization"
)

type JWTProvider interface {
	ValidateToken(token string) error
}

type Middleware struct {
	provider JWTProvider
}

func NewMiddleware(provider JWTProvider) *Middleware {
	m := &Middleware{
		provider: provider,
	}

	return m
}

func (m *Middleware) AccessTokenValidation() gin.HandlerFunc {
	return func(c *gin.Context) {

		accessToken, err := ExtractTokenFromHeader(c)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.BuildJSONResponse("Error", err.Error()))
			return
		}

		err = m.provider.ValidateToken(accessToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.BuildJSONResponse("Error", err.Error()))
			return
		}

		login, err := ExtractLoginFromToken(accessToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.BuildJSONResponse("Error", err.Error()))
			return
		}

		c.Set("login", login)
		c.Next()
	}
}

func ExtractTokenFromHeader(c *gin.Context) (string, error) {
	rawAccessToken := c.GetHeader(Header)
	if rawAccessToken == "" {
		return "", errors.New("no access token in headers")
	}
	parts := strings.Split(rawAccessToken, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("invalid token format")
	}

	return parts[1], nil
}

func ExtractLoginFromToken(tokenStr string) (string, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenStr, jwt.MapClaims{})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("invalid token payload")
	}

	login, ok := claims["login"].(string)

	if !ok {
		return "", fmt.Errorf("invalid token payload")
	}
	return login, nil
}
