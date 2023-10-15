package handler

import (
	"GatewayService/internal/handler/error/validator"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"net/http"
)

type StoresHandler struct {
	logger          *zap.Logger
	rabbitMQChannel *amqp.Channel
	rabbitMQConn    *amqp.Connection
	rabbitMQQueue   string
}

type Store struct {
	Name        string `json:"name" binding:"required"`
	Address     string `json:"address" binding:"required"`
	OwnerName   string `json:"ownerName" binding:"required"`
	OpeningTime string `json:"openingTime" binding:"required"`
	ClosingTime string `json:"closingTime" binding:"required"`
}

type StoreVersion struct {
	StoreOwnerName string `json:"storeOwnerName" binding:"required"`
	OpeningTime    string `json:"openingTime" binding:"required"`
	ClosingTime    string `json:"closingTime" binding:"required"`
}

func NewStoresHandler(channel *amqp.Channel, rabbitMQConn *amqp.Connection, rabbitMQQueue string, logger *zap.Logger) *StoresHandler {
	return &StoresHandler{
		logger:          logger,
		rabbitMQChannel: channel,
		rabbitMQConn:    rabbitMQConn,
		rabbitMQQueue:   rabbitMQQueue,
	}
}

func (h *StoresHandler) sendMessage(message []byte) error {
	err := h.rabbitMQChannel.Publish(
		"",
		h.rabbitMQQueue,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func buildMessage(data interface{}, action, login, storeId, versionId string) []byte {
	message := map[string]interface{}{
		"storeId":   storeId,
		"versionId": versionId,
		"data":      data,
		"action":    action,
		"userLogin": login,
	}

	body, err := json.Marshal(message)
	if err != nil {
		return nil
	}

	return body
}

func (h *StoresHandler) CreateStore(c *gin.Context) {
	var store Store
	if err := c.ShouldBindJSON(&store); err != nil {
		c.JSON(http.StatusBadRequest, validator.ProcessValidatorError(err))
		return
	}

	action := "create_store"

	login := c.GetString("login")

	err := h.sendMessage(buildMessage(store, action, login, "", ""))

	if err != nil {
		h.logger.With(
			zap.String("place", "Handler"),
			zap.Error(err),
		).Error("Failed to publish a message")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish a message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Message sent to storage"})
}

func (h *StoresHandler) CreateStoreVersion(c *gin.Context) {
	var storeVersion StoreVersion
	if err := c.ShouldBindJSON(&storeVersion); err != nil {
		c.JSON(http.StatusBadRequest, validator.ProcessValidatorError(err))
		return
	}

	action := "create_store_version"

	login := c.GetString("login")

	storeId := c.Param("id")

	err := h.sendMessage(buildMessage(storeVersion, action, login, storeId, ""))

	if err != nil {
		h.logger.With(
			zap.String("place", "Handler"),
			zap.Error(err),
		).Error("Failed to publish a message")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish a message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Message sent to storage"})
}

func (h *StoresHandler) DeleteStore(c *gin.Context) {
	action := "delete_store"

	login := c.GetString("login")

	storeId := c.Param("id")

	err := h.sendMessage(buildMessage(nil, action, login, storeId, ""))

	if err != nil {
		h.logger.With(
			zap.String("place", "Handler"),
			zap.Error(err),
		).Error("Failed to publish a message")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish a message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Message sent to storage"})
}

func (h *StoresHandler) DeleteStoreVersion(c *gin.Context) {
	action := "delete_store_version"

	login := c.GetString("login")

	storeId := c.Param("id")

	versionId := c.Param("versionId")

	err := h.sendMessage(buildMessage(nil, action, login, storeId, versionId))

	if err != nil {
		h.logger.With(
			zap.String("place", "Handler"),
			zap.Error(err),
		).Error("Failed to publish a message")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish a message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Message sent to storage"})
}

func (h *StoresHandler) GetStore(c *gin.Context) {
	action := "get_store"

	login := c.GetString("login")

	storeId := c.Param("id")

	err := h.sendMessage(buildMessage(nil, action, login, storeId, ""))

	if err != nil {
		h.logger.With(
			zap.String("place", "Handler"),
			zap.Error(err),
		).Error("Failed to publish a message")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish a message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Message sent to storage"})
}

func (h *StoresHandler) GetStoreHistory(c *gin.Context) {

	action := "get_store_history"

	login := c.GetString("login")

	storeId := c.Param("id")

	err := h.sendMessage(buildMessage(nil, action, login, storeId, ""))
	if err != nil {
		h.logger.With(
			zap.String("place", "Handler"),
			zap.Error(err),
		).Error("Failed to publish a message")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish a message"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Message sent to storage"})
}

func (h *StoresHandler) GetStoreVersion(c *gin.Context) {

	action := "get_store_version"

	login := c.GetString("login")

	storeId := c.Param("id")

	versionId := c.Param("versionId")

	err := h.sendMessage(buildMessage(nil, action, login, storeId, versionId))
	if err != nil {
		h.logger.With(
			zap.String("place", "Handler"),
			zap.Error(err),
		).Error("Failed to publish a message")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish a message"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Message sent to storage"})
}
