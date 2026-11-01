package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SpecElectricalPanelHandler struct {
	DB *gorm.DB
}

func (h *SpecElectricalPanelHandler) GetSpecElectricalPanels(c *gin.Context) {
	var specElectricalPanels []SpecElectricalPanel
	if err := h.DB.Find(&specElectricalPanels).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, specElectricalPanels)
}

func (h *SpecElectricalPanelHandler) GetSpecElectricalPanel(c *gin.Context) {
	id := c.Param("spec_electrical_panel_id")
	var spec SpecElectricalPanel
	if err := h.DB.First(&spec, "equipment_id = ?", id).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, spec)
}

func (h *SpecElectricalPanelHandler) CreateSpecElectricalPanel(c *gin.Context) {
	var spec SpecElectricalPanel
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

func (h *SpecElectricalPanelHandler) UpdateSpecElectricalPanel(c *gin.Context) {
	id := c.Param("spec_electrical_panel_id")
	var spec SpecElectricalPanel
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

func (h *SpecElectricalPanelHandler) DeleteSpecElectricalPanel(c *gin.Context) {
	id := c.Param("spec_electrical_panel_id")
	if err := h.DB.Delete(&SpecElectricalPanel{}, "equipment_id = ?", id).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.Status(204)
}
