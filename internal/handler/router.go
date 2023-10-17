package handler

import (
	"GatewayService/internal/middleware"
	"github.com/gin-gonic/gin"
)

func NewRouter(authHandler *AuthHandler, storesHandler *StoresHandler, middleware *middleware.Middleware) *gin.Engine {
	router := gin.Default()

	authGroup := router.Group("auth")
	authGroup.POST("/login", authHandler.SingIn)

	storesGroup := router.Group("storage")
	storesGroup.POST("/store", middleware.AccessTokenValidation(), storesHandler.CreateStore)
	storesGroup.POST("/store/:id/version", middleware.AccessTokenValidation(), storesHandler.CreateStoreVersion)
	storesGroup.DELETE("/store/:id", middleware.AccessTokenValidation(), storesHandler.DeleteStore)
	storesGroup.DELETE("/store/:id/version/:versionId", middleware.AccessTokenValidation(), storesHandler.DeleteStoreVersion)
	storesGroup.GET("/store/:id", middleware.AccessTokenValidation(), storesHandler.GetStore)
	storesGroup.GET("/store/:id/history", middleware.AccessTokenValidation(), storesHandler.GetStoreHistory)
	storesGroup.GET("/store/:id/version/:versionId", middleware.AccessTokenValidation(), storesHandler.GetStoreVersion)

	//for response handling from storage service
	responseGroup := router.Group("response")
	responseGroup.POST("/", storesHandler.HandleResponse)

	return router
}
