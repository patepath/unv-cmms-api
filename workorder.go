package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type WorkOrderHandler struct {
	DB *gorm.DB
}

func (h *WorkOrderHandler) GetWorkOrders(c *gin.Context) {
	var workOrders []WorkOrder
	if err := h.DB.Order("scheduled_date, time").Find(&workOrders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, workOrders)
}

func (h *WorkOrderHandler) CreateWorkOrder(c *gin.Context) {
	var workOrder WorkOrder
	if err := c.ShouldBindJSON(&workOrder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if workOrder.Title == "" || workOrder.Asset == "" || workOrder.ScheduledDate == "" || workOrder.Time == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title, asset, scheduledDate and time are required"})
		return
	}
	if err := h.DB.Create(&workOrder).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, workOrder)
}

func (h *WorkOrderHandler) UpdateWorkOrder(c *gin.Context) {
	var workOrder WorkOrder
	if err := h.DB.First(&workOrder, c.Param("id")).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Work order not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var input struct {
		Title         *string `json:"title"`
		Asset         *string `json:"asset"`
		Priority      *string `json:"priority"`
		Assignee      *string `json:"assignee"`
		Status        *string `json:"status"`
		ScheduledDate *string `json:"scheduledDate"`
		Time          *string `json:"time"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updates := map[string]any{}
	if input.Title != nil {
		updates["title"] = *input.Title
	}
	if input.Asset != nil {
		updates["asset"] = *input.Asset
	}
	if input.Priority != nil {
		updates["priority"] = *input.Priority
	}
	if input.Assignee != nil {
		updates["assignee"] = *input.Assignee
	}
	if input.Status != nil {
		updates["status"] = *input.Status
	}
	if input.ScheduledDate != nil {
		updates["scheduled_date"] = *input.ScheduledDate
	}
	if input.Time != nil {
		updates["time"] = *input.Time
	}
	if err := h.DB.Model(&workOrder).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	h.DB.First(&workOrder, workOrder.ID)
	c.JSON(http.StatusOK, workOrder)
}
