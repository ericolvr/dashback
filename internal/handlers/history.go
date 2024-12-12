package handlers

import (
	"context"
	"net/http"

	"github.com/Alarmtekgit/websocket/internal/domain"
	"github.com/Alarmtekgit/websocket/internal/service"
	"github.com/gin-gonic/gin"
)

type HistoryHandler struct {
	service service.HistoryService
}

func NewHistoryHandler(service service.HistoryService) *HistoryHandler {
	return &HistoryHandler{service: service}
}

func (h *HistoryHandler) CreateHistory(c *gin.Context) {
	var history domain.History
	if err := c.ShouldBindJSON(&history); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := h.service.CreateHistory(context.Background(), &history); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create message"})
		return
	}

	c.JSON(http.StatusCreated, history)
}

func (h *HistoryHandler) FindHistoryByID(c *gin.Context) {
	equipmentid := c.Param("equipmentid")
	history, err := h.service.FindHistoryByID(context.Background(), equipmentid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "History not found"})
		return
	}

	c.JSON(http.StatusOK, history)
}

func (h *HistoryHandler) GetHistoryByID(c *gin.Context) {
	equipmentid := c.Param("equipmentid")
	histories, err := h.service.GetHistoryByID(context.Background(), equipmentid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch histories"})
		return
	}

	c.JSON(http.StatusOK, histories)
}
