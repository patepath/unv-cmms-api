package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SpecTransformerHandler struct {
	DB *gorm.DB
}

func (h *SpecTransformerHandler) GetSpecTransformers(c *gin.Context) {
	var specTransformers []SpecTransformer
	if err := h.DB.Find(&specTransformers).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to get spec transformers"})
		return
	}

	c.JSON(200, specTransformers)
}

func (h *SpecTransformerHandler) GetSpecTransformer(c *gin.Context) {
	specTransformerID := c.Param("spec_transformer_id")

	var specTransformer SpecTransformer
	if err := h.DB.First(&specTransformer, "spec_transformer_id = ?", specTransformerID).Error; err != nil {
		c.JSON(404, gin.H{"error": "Spec transformer not found"})
		return
	}

	c.JSON(200, specTransformer)
}

func (h *SpecTransformerHandler) CreateSpecTransformer(c *gin.Context) {
	var specTransformer SpecTransformer
	if err := c.ShouldBindJSON(&specTransformer); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	if err := h.DB.Create(&specTransformer).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to create spec transformer"})
		return
	}

	c.JSON(201, specTransformer)
}

func (h *SpecTransformerHandler) UpdateSpecTransformer(c *gin.Context) {
	specTransformerID := c.Param("spec_transformer_id")

	var specTransformer SpecTransformer
	if err := h.DB.First(&specTransformer, "spec_transformer_id = ?", specTransformerID).Error; err != nil {
		c.JSON(404, gin.H{"error": "Spec transformer not found"})
		return
	}

	if err := c.ShouldBindJSON(&specTransformer); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	if err := h.DB.Save(&specTransformer).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to update spec transformer"})
		return
	}

	c.JSON(200, specTransformer)
}

func (h *SpecTransformerHandler) DeleteSpecTransformer(c *gin.Context) {
	specTransformerID := c.Param("spec_transformer_id")

	if err := h.DB.Delete(&SpecTransformer{}, "spec_transformer_id = ?", specTransformerID).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete spec transformer"})
		return
	}

	c.JSON(200, gin.H{"message": "Spec transformer deleted successfully"})
}
