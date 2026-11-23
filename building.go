package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BuildingHandler struct {
	DB *gorm.DB
}

func (h *BuildingHandler) GetBuildings(c *gin.Context) {
	var buildings []Building
	if err := h.DB.Order("code").Find(&buildings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, buildings)
}

func (h *BuildingHandler) GetBuilding(c *gin.Context) {
	buildingCode := c.Param("building_code")
	var building Building
	if err := h.DB.Where("code = ?", buildingCode).First(&building).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Building not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, building)
}

func (h *BuildingHandler) CreateBuilding(c *gin.Context) {
	var building Building
	if err := c.ShouldBindJSON(&building); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.DB.Create(&building).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, building)
}

func (h *BuildingHandler) UpdateBuilding(c *gin.Context) {
	buildingCode := c.Param("building_code")
	var building Building
	if err := h.DB.Where("code = ?", buildingCode).First(&building).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Building not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := c.ShouldBindJSON(&building); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.DB.Save(&building).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, building)
}

func (h *BuildingHandler) DeleteBuilding(c *gin.Context) {
	buildingCode := c.Param("building_code")
	if err := h.DB.Where("code = ?", buildingCode).Delete(&Building{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Building deleted successfully"})
}
