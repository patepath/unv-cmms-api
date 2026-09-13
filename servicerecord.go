package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ServiceRecordHandler struct {
	DB *gorm.DB
}

func (h *ServiceRecordHandler) GetServiceRecords(c *gin.Context) {
	var serviceRecords []ServiceRecord
	h.DB.Order("service_date").Find(&serviceRecords)
	c.JSON(http.StatusOK, serviceRecords)
}

func (h *ServiceRecordHandler) GetServiceRecord(c *gin.Context) {
	var serviceRecord ServiceRecord
	id := c.Param("id")

	if err := h.DB.First(&serviceRecord, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Service record not found"})
		return
	}

	c.JSON(http.StatusOK, serviceRecord)
}
func (h *ServiceRecordHandler) CreateServiceRecord(c *gin.Context) {
	var serviceRecord ServiceRecord
	if err := c.ShouldBindJSON(&serviceRecord); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.DB.Create(&serviceRecord).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, serviceRecord)
}

func (h *ServiceRecordHandler) UpdateServiceRecord(c *gin.Context) {
	var serviceRecord ServiceRecord
	id := c.Param("id")

	if err := h.DB.First(&serviceRecord, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Service record not found"})
		return
	}

	if err := c.ShouldBindJSON(&serviceRecord); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.DB.Save(&serviceRecord).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, serviceRecord)
}

func (h *ServiceRecordHandler) DeleteServiceRecord(c *gin.Context) {
	var serviceRecord ServiceRecord
	id := c.Param("id")

	if err := h.DB.First(&serviceRecord, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Service record not found"})
		return
	}

	if err := h.DB.Delete(&serviceRecord).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Service record deleted successfully"})
}
