package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SpecPumpHandler struct {
	DB *gorm.DB
}

func (h *SpecPumpHandler) GetSpecPumps(c *gin.Context) {
	var specPumps []SpecPump
	if err := h.DB.Find(&specPumps).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, specPumps)
}

func (h *SpecPumpHandler) GetSpecPump(c *gin.Context) {
	id := c.Param("spec_pump_id")
	var spec SpecPump
	if err := h.DB.First(&spec, "equipment_id = ?", id).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, spec)
}

func (h *SpecPumpHandler) CreateSpecPump(c *gin.Context) {
	var spec SpecPump
	if err := c.ShouldBindJSON(&spec); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := h.DB.Create(&spec).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, spec)
}

func (h *SpecPumpHandler) UpdateSpecPump(c *gin.Context) {
	id := c.Param("spec_pump_id")
	var spec SpecPump
	if err := h.DB.First(&spec, "equipment_id = ?", id).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if err := c.ShouldBindJSON(&spec); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := h.DB.Save(&spec).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, spec)
}

func (h *SpecPumpHandler) DeleteSpecPump(c *gin.Context) {
	id := c.Param("spec_pump_id")
	var spec SpecPump
	if err := h.DB.First(&spec, "equipment_id = ?", id).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if err := h.DB.Delete(&spec).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "SpecPump deleted successfully"})
}
