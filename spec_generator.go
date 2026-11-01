package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SpecGeneratorHandler struct {
	DB *gorm.DB
}

func (h *SpecGeneratorHandler) GetSpecGenerators(c *gin.Context) {
	var specGenerators []SpecGenerator
	if err := h.DB.Find(&specGenerators).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to get spec generators"})
		return
	}

	c.JSON(200, specGenerators)
}

func (h *SpecGeneratorHandler) GetSpecGenerator(c *gin.Context) {
	specGeneratorID := c.Param("spec_generator_id")

	var specGenerator SpecGenerator
	if err := h.DB.First(&specGenerator, "equipment_id = ?", specGeneratorID).Error; err != nil {
		c.JSON(404, gin.H{"error": "Spec generator not found"})
		return
	}

	c.JSON(200, specGenerator)
}

func (h *SpecGeneratorHandler) CreateSpecGenerator(c *gin.Context) {
	var specGenerator SpecGenerator
	if err := c.ShouldBindJSON(&specGenerator); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	if err := h.DB.Create(&specGenerator).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to create spec generator"})
		return
	}

	c.JSON(201, specGenerator)
}
func (h *SpecGeneratorHandler) UpdateSpecGenerator(c *gin.Context) {
	specGeneratorID := c.Param("spec_generator_id")

	var specGenerator SpecGenerator
	if err := h.DB.First(&specGenerator, "equipment_id = ?", specGeneratorID).Error; err != nil {
		c.JSON(404, gin.H{"error": "Spec generator not found"})
		return
	}

	if err := c.ShouldBindJSON(&specGenerator); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	if err := h.DB.Save(&specGenerator).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to update spec generator"})
		return
	}

	c.JSON(200, specGenerator)
}

func (h *SpecGeneratorHandler) DeleteSpecGenerator(c *gin.Context) {
	specGeneratorID := c.Param("spec_generator_id")

	var specGenerator SpecGenerator
	if err := h.DB.First(&specGenerator, "equipment_id = ?", specGeneratorID).Error; err != nil {
		c.JSON(404, gin.H{"error": "Spec generator not found"})
		return
	}

	if err := h.DB.Delete(&specGenerator).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete spec generator"})
		return
	}

	c.JSON(200, gin.H{"message": "Spec generator deleted"})
}
