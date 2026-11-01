package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BuildingHandler struct {
	DB *gorm.DB
}

func (h *BuildingHandler) GetBuildings(c *gin.Context) {
	var buildings []Building
	if err := h.DB.Order("code").Find(&buildings).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, buildings)
}

func (h *BuildingHandler) GetBuilding(c *gin.Context) {
	buildingCode := c.Param("building_code")
	var building Building
	if err := h.DB.Where("building_code = ?", buildingCode).First(&building).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, building)
}

func (h *BuildingHandler) CreateBuilding(c *gin.Context) {
	var building Building
	if err := c.ShouldBindJSON(&building); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := h.DB.Create(&building).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, building)
}

func (h *BuildingHandler) UpdateBuilding(c *gin.Context) {
	buildingCode := c.Param("building_code")
	var building Building
	if err := h.DB.Where("building_code = ?", buildingCode).First(&building).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if err := c.ShouldBindJSON(&building); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := h.DB.Save(&building).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, building)
}

func (h *BuildingHandler) DeleteBuilding(c *gin.Context) {
	buildingCode := c.Param("building_code")
	if err := h.DB.Where("building_code = ?", buildingCode).Delete(&Building{}).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Building deleted successfully"})
}
