package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"net/http"
)

type StoresHandler struct {
	logger  *zap.Logger
	channel *amqp.Channel
}

type Store struct {
	Name    string `json:"name" binding:"required,min=3,max=40"`
	Address string `json:"address" binding:"required,min=6,max=40"  minimum:"6" maximum:"40" default:"password"`
}

func NewStoresHandler(channel *amqp.Channel, logger *zap.Logger) *StoresHandler {
	return &StoresHandler{
		logger:  logger,
		channel: channel,
	}
}

func (h *StoresHandler) CreateStore(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Entity created"})
}

func (h *StoresHandler) CreateStoreVersion(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Entity version created"})
}

func (h *StoresHandler) DeleteStore(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Entity deleted"})
}

func (h *StoresHandler) DeleteStoreVersion(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Entity version deleted"})
}

func (h *StoresHandler) GetStore(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Get entity"})
}

func (h *StoresHandler) GetStoreHistory(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Get entity history"})
}

func (h *StoresHandler) GetStoreVersion(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Get entity version"})
}
