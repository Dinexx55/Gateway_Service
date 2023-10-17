package handler

import (
	"GatewayService/internal/handler/mapper"
	"GatewayService/internal/handler/response"
	"GatewayService/internal/handler/validation"
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
	Login    string `json:"login" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6,max=40"`
}

func NewAuthHandler(authService AuthService, logger *zap.Logger, mapper mapper.ErrorMapper) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		logger:      logger,
		errorMapper: mapper,
	}
}

func (h *AuthHandler) SingIn(c *gin.Context) {

	var credentials Auth
	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, validation.FormatValidatorError(err))
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
		).Error("Error while signing in: " + err.Error())

		errInf := h.errorMapper.MapError(err)

		c.JSON(errInf.StatusCode,
			response.BuildJSONResponse("Error", errInf.Message))

		return
	}

	h.logger.With(
		zap.String("token", "accessToken"),
	).Info("Token generated successfully")

	c.JSON(http.StatusOK, response.BuildJSONResponse("Access token", accessToken))
}
