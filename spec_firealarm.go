package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SpecFireAlarmHandler struct {
	DB *gorm.DB
}

func (h *SpecFireAlarmHandler) GetSpecFireAlarms(c *gin.Context) {
	var specFireAlarms []SpecFireAlarm
	if err := h.DB.Find(&specFireAlarms).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, specFireAlarms)
}

func (h *SpecFireAlarmHandler) GetSpecFireAlarm(c *gin.Context) {
	id := c.Param("spec_fire_alarm_id")
	var spec SpecFireAlarm
	if err := h.DB.First(&spec, "equipment_id = ?", id).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, spec)
}

func (h *SpecFireAlarmHandler) CreateSpecFireAlarm(c *gin.Context) {
	var spec SpecFireAlarm
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

func (h *SpecFireAlarmHandler) UpdateSpecFireAlarm(c *gin.Context) {
	id := c.Param("spec_fire_alarm_id")
	var spec SpecFireAlarm
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

func (h *SpecFireAlarmHandler) DeleteSpecFireAlarm(c *gin.Context) {
	id := c.Param("spec_fire_alarm_id")
	if err := h.DB.Delete(&SpecFireAlarm{}, "equipment_id = ?", id).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.Status(204)
}
