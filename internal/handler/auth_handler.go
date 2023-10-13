package handler

import (
	"GatewayService/internal/handler/error/mapper"
	"GatewayService/internal/handler/error/validator"
	"GatewayService/internal/handler/response"
	"GatewayService/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type AuthService interface {
	SignIn(user service.User) (string, error)
}

type AuthHandler struct {
	authService AuthService
	logger      *zap.Logger
	errorMapper mapper.ErrorMapper
}

type Auth struct {
	Login    string `json:"login" binding:"required,min=3,max=40"`
	Password string `json:"password" binding:"required,min=6,max=40"  minimum:"6" maximum:"40" default:"password"`
}

func NewAuthHandler(authService AuthService, logger *zap.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		logger:      logger,
	}
}

func (h *AuthHandler) SingIn(c *gin.Context) {

	var credentials Auth
	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, validator.ProcessValidatorError(err))
		return
	}

	user := service.User{
		Login:    credentials.Login,
		Password: credentials.Password,
	}

	accessToken, err := h.authService.SignIn(user)

	if err != nil {
		h.logger.With(
			zap.String("place", "authHandler"),
			zap.String("func", "SignIn"),
		).Error("Error while signing in")

		errInf := h.errorMapper.MapError(err)
		c.JSON(errInf.StatusCode,
			response.CreateJSONResult("Error", errInf.Msg))

		return
	}

	h.logger.With(
		zap.String("token", "accessToken"),
	).Info("Generated successfully")

	c.JSON(http.StatusOK, response.CreateJSONResult("Access token", accessToken))
}