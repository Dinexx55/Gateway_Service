package handler

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ShopsHandler struct {
	logger *zap.Logger
}

func (h ShopsHandler) CreateShop(context *gin.Context) {

}

func (h ShopsHandler) GetShopById(context *gin.Context) {

}

func NewShopsHandler(logger *zap.Logger) *ShopsHandler {
	return &ShopsHandler{
		logger: logger,
	}
}
