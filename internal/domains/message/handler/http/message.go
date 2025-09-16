package http

import (
	"net/http"
	"time"

	"MessageService/internal/domains/message/usecase"

	"github.com/gin-gonic/gin"
)

type MessageHandler struct {
	uc usecase.MessageUseCase
}

func NewMessageHandler(uc usecase.MessageUseCase) *MessageHandler {
	return &MessageHandler{uc: uc}
}

func (h *MessageHandler) RegisterRoutes(router *gin.Engine) {
	router.GET("/messages/:id", h.GetByID)
	router.GET("/messages", h.GetByTimeRange)
}

func (h *MessageHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	msg, err := h.uc.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "message not found", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, msg)
}

func (h *MessageHandler) GetByTimeRange(c *gin.Context) {
	fromStr := c.Query("from")
	toStr := c.Query("to")

	if fromStr == "" || toStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing 'from' or 'to' query params"})
		return
	}

	from, err := time.Parse(time.RFC3339, fromStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid 'from' format, use RFC3339"})
		return
	}

	to, err := time.Parse(time.RFC3339, toStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid 'to' format, use RFC3339"})
		return
	}

	messages, err := h.uc.GetByTimeRange(c.Request.Context(), from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get messages", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, messages)
}
