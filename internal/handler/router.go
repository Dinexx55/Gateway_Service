package handler

import (
	"GatewayService/internal/middleware"
	"github.com/gin-gonic/gin"
)

func NewRouter(authHandler *AuthHandler, shopsHandler *ShopsHandler, middleware *middleware.Middleware) *gin.Engine {
	router := gin.Default()

	authGroup := router.Group("auth")
	{
		authGroup.POST("/login", authHandler.SingIn)
	}

	storageGroup := router.Group("storage")
	{
		storageGroup.POST("/shop", middleware.AccessTokenValidation(), shopsHandler.CreateShop)
		storageGroup.GET("/shop/:shopId", middleware.AccessTokenValidation(), shopsHandler.GetShopById)
	}
	return router
}
