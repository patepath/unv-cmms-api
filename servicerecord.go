package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ServiceRecordHandler struct {
	DB *gorm.DB
}

type repeatedServiceRecordRequest struct {
	AssetID        uint      `json:"assetId"`
	RecordType     string    `json:"recordType"`
	ServiceDate    time.Time `json:"serviceDate"`
	Details        string    `json:"details"`
	RepeatInterval string    `json:"repeatInterval"`
	RepeatCount    int       `json:"repeatCount"`
}

func (h *ServiceRecordHandler) GetServiceRecords(c *gin.Context) {
	var serviceRecords []ServiceRecord
	query := h.DB.Order("service_date")
	if recordType := c.Query("type"); recordType != "" {
		query = query.Where("record_type = ?", recordType)
	}
	if err := query.Find(&serviceRecords).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
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
	if serviceRecord.EndDate != nil && !serviceRecord.EndDate.After(serviceRecord.ServiceDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "endDate must be after serviceDate"})
		return
	}

	if err := h.DB.Create(&serviceRecord).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, serviceRecord)
}

func (h *ServiceRecordHandler) CreateRepeatedServiceRecords(c *gin.Context) {
	var request repeatedServiceRecordRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if request.AssetID == 0 || request.Details == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "assetId and details are required"})
		return
	}
	if request.RepeatCount < 1 || request.RepeatCount > 366 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "repeatCount must be between 1 and 366"})
		return
	}

	records := make([]ServiceRecord, 0, request.RepeatCount)
	for index := 0; index < request.RepeatCount; index++ {
		serviceDate, ok := repeatServiceDate(request.ServiceDate, request.RepeatInterval, index)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid repeat interval"})
			return
		}
		records = append(records, ServiceRecord{
			AssetID:     request.AssetID,
			RecordType:  request.RecordType,
			ServiceDate: serviceDate,
			Details:     request.Details,
		})
	}

	if err := h.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&records).Error
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, records)
}

func repeatServiceDate(start time.Time, interval string, occurrence int) (time.Time, bool) {
	switch interval {
	case "daily":
		return start.AddDate(0, 0, occurrence), true
	case "weekly":
		return start.AddDate(0, 0, occurrence*7), true
	case "every15days":
		return start.AddDate(0, 0, occurrence*15), true
	case "monthly":
		return start.AddDate(0, occurrence, 0), true
	case "quarterly":
		return start.AddDate(0, occurrence*3, 0), true
	case "halfYearly":
		return start.AddDate(0, occurrence*6, 0), true
	case "yearly":
		return start.AddDate(occurrence, 0, 0), true
	default:
		return time.Time{}, false
	}
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
