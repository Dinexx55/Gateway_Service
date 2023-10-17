package handler

import (
	"GatewayService/internal/handler/response"
	"GatewayService/internal/handler/validation"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"net/http"
)

type StoresHandler struct {
	logger          *zap.Logger
	rabbitMQChannel *amqp.Channel
	rabbitMQConn    *amqp.Connection
	rabbitMQQueue   string
	structValidator *validator.Validate
}

// Some custom validators used
type Store struct {
	Name        string `json:"name" validate:"required,min=3,max=40"`
	Address     string `json:"address" validate:"required,addressFormat"`
	OwnerName   string `json:"ownerName" validate:"required,ownerNameFormat"`
	OpeningTime string `json:"openingTime" validate:"required,timeFormat"`
	ClosingTime string `json:"closingTime" validate:"required,timeFormat"`
}

type StoreVersion struct {
	OwnerName   string `json:"ownerName" validate:"required,ownerNameFormat"`
	OpeningTime string `json:"openingTime" validate:"required,timeFormat"`
	ClosingTime string `json:"closingTime" validate:"required,timeFormat"`
}

func NewStoresHandler(channel *amqp.Channel, rabbitMQConn *amqp.Connection, rabbitMQQueue string, logger *zap.Logger, structValidator *validator.Validate) *StoresHandler {
	return &StoresHandler{
		logger:          logger,
		rabbitMQChannel: channel,
		rabbitMQConn:    rabbitMQConn,
		rabbitMQQueue:   rabbitMQQueue,
		structValidator: structValidator,
	}
}

const (
	messageForError   = "Failed to publish a message"
	messageForSuccess = "Storage service is processing your message. Check status through logs"
)

func (h *StoresHandler) CreateStore(c *gin.Context) {
	var store Store
	if err := c.ShouldBindJSON(&store); err != nil {
		c.JSON(http.StatusBadRequest, validation.FormatValidatorError(err))
		return
	}

	if err := h.structValidator.Struct(store); err != nil {
		c.JSON(http.StatusBadRequest, validation.FormatValidatorError(err))
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
		c.JSON(http.StatusInternalServerError, response.BuildJSONResponse("Error", messageForError))
		return
	}

	c.JSON(http.StatusOK, response.BuildJSONResponse("Success", messageForSuccess))
}

func (h *StoresHandler) CreateStoreVersion(c *gin.Context) {
	var storeVersion StoreVersion
	if err := c.ShouldBindJSON(&storeVersion); err != nil {
		c.JSON(http.StatusBadRequest, validation.FormatValidatorError(err))
		return
	}

	if err := h.structValidator.Struct(storeVersion); err != nil {
		c.JSON(http.StatusBadRequest, validation.FormatValidatorError(err))
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
		c.JSON(http.StatusInternalServerError, response.BuildJSONResponse("Error", messageForError))
		return
	}

	c.JSON(http.StatusOK, response.BuildJSONResponse("Success", messageForSuccess))
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
		c.JSON(http.StatusInternalServerError, response.BuildJSONResponse("Error", messageForError))
		return
	}

	c.JSON(http.StatusOK, response.BuildJSONResponse("Success", messageForSuccess))
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
		c.JSON(http.StatusInternalServerError, response.BuildJSONResponse("Error", messageForError))
		return
	}

	c.JSON(http.StatusOK, response.BuildJSONResponse("Success", messageForSuccess))
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
		c.JSON(http.StatusInternalServerError, response.BuildJSONResponse("Error", messageForError))
		return
	}

	c.JSON(http.StatusOK, response.BuildJSONResponse("Success", messageForSuccess))
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
		c.JSON(http.StatusInternalServerError, response.BuildJSONResponse("Error", messageForError))
		return
	}
	c.JSON(http.StatusOK, response.BuildJSONResponse("Success", messageForSuccess))
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
		c.JSON(http.StatusInternalServerError, response.BuildJSONResponse("Error", messageForError))
		return
	}
	c.JSON(http.StatusOK, response.BuildJSONResponse("Success", messageForSuccess))
}

func (h *StoresHandler) HandleResponse(c *gin.Context) {
	var payload interface{}

	h.logger.Info("Trying to extract payload from response from storage service")

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Printf("Received payload: %+v\n", payload)

	c.JSON(http.StatusOK, payload)
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
