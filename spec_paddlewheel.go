package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SpecPaddleWheelHandler struct {
	DB *gorm.DB
}

func (h *SpecPaddleWheelHandler) GetSpecPaddleWheels(c *gin.Context) {
	var specPaddleWheels []SpecPaddleWheel
	if err := h.DB.Find(&specPaddleWheels).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, specPaddleWheels)
}

func (h *SpecPaddleWheelHandler) GetSpecPaddleWheel(c *gin.Context) {
	specPaddleWheelID := c.Param("spec_paddle_wheel_id")
	var specPaddleWheel SpecPaddleWheel
	if err := h.DB.First(&specPaddleWheel, "equipment_id = ?", specPaddleWheelID).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, specPaddleWheel)
}

func (h *SpecPaddleWheelHandler) CreateSpecPaddleWheel(c *gin.Context) {
	var specPaddleWheel SpecPaddleWheel
	if err := c.ShouldBindJSON(&specPaddleWheel); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := h.DB.Create(&specPaddleWheel).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, specPaddleWheel)
}

func (h *SpecPaddleWheelHandler) UpdateSpecPaddleWheel(c *gin.Context) {
	specPaddleWheelID := c.Param("spec_paddle_wheel_id")
	var specPaddleWheel SpecPaddleWheel
	if err := h.DB.First(&specPaddleWheel, "equipment_id = ?", specPaddleWheelID).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if err := c.ShouldBindJSON(&specPaddleWheel); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := h.DB.Save(&specPaddleWheel).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, specPaddleWheel)
}

func (h *SpecPaddleWheelHandler) DeleteSpecPaddleWheel(c *gin.Context) {
	specPaddleWheelID := c.Param("spec_paddle_wheel_id")
	if err := h.DB.Delete(&SpecPaddleWheel{}, "equipment_id = ?", specPaddleWheelID).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "SpecPaddleWheel deleted successfully"})
}
