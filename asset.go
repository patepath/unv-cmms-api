package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AssetHandler struct {
	DB *gorm.DB
}

func (h *AssetHandler) GetAssets(c *gin.Context) {
	var assets []Asset

	if err := h.DB.Preload("SpecTransformer").
		Preload("SpecGenerator").
		Preload("SpecElectricalPanel").
		Preload("SpecPump").
		Preload("SpecPaddleWheel").
		Find(&assets).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch assets"})
		return
	}

	c.JSON(200, assets)
}

func (h *AssetHandler) GetAsset(c *gin.Context) {
	assetID := c.Param("asset_id")

	var asset Asset
	if err := h.DB.Preload("SpecTransformer").
		Preload("SpecGenerator").
		Preload("SpecElectricalPanel").
		Preload("SpecPump").
		Preload("SpecPaddleWheel").
		First(&asset, "id = ?", assetID).Error; err != nil {
		c.JSON(404, gin.H{"error": "Asset not found"})
		return
	}

	c.JSON(200, asset)
}

func (h *AssetHandler) CreateAsset(c *gin.Context) {
	var asset Asset
	if err := c.ShouldBindJSON(&asset); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	if err := h.DB.Create(&asset).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to create asset"})
		return
	}

	c.JSON(201, asset)
}

func (h *AssetHandler) UpdateAsset(c *gin.Context) {
	assetID := c.Param("id")

	var asset Asset
	if err := h.DB.First(&asset, "id = ?", assetID).Error; err != nil {
		c.JSON(404, gin.H{"error": "Asset not found"})
		return
	}

	if err := c.ShouldBindJSON(&asset); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	if err := h.DB.Save(&asset).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to update asset"})
		return
	}

	c.JSON(200, asset)
}

func (h *AssetHandler) DeleteAsset(c *gin.Context) {
	assetID := c.Param("id")

	if err := h.DB.Delete(&Asset{}, "id = ?", assetID).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete asset"})
		return
	}

	c.JSON(200, gin.H{"message": "Asset deleted successfully"})
}
