package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AssetHandler struct {
	DB *gorm.DB
}

func (h *AssetHandler) GetAssetsByCategory(c *gin.Context) {
	categoryID := c.Param("category_id")
	var assets []Asset

	if err := h.DB.Preload("Category").
		Preload("Brand").
		Preload("Buildings").
		Preload("SpecTransformer").
		Preload("SpecGenerator").
		Preload("SpecElectricalPanel").
		Preload("SpecFireAlarm").
		Preload("SpecPump").
		Preload("SpecPaddleWheel").
		Where("category_id = ?", categoryID).
		Find(&assets).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch assets by category"})
		return
	}

	c.JSON(200, assets)
}

func (h *AssetHandler) GetAssets(c *gin.Context) {
	var assets []Asset

	if err := h.DB.Preload("SpecTransformer").
		Preload("SpecGenerator").
		Preload("SpecElectricalPanel").
		Preload("SpecPump").
		Preload("SpecPaddleWheel").
		Preload("Buildings").
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
		Preload("Buildings").
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

	buildings, err := h.getBuildings(asset.BuildingIDs)
	if err != nil {
		c.JSON(400, gin.H{"error": "One or more buildings were not found"})
		return
	}

	if err := h.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Omit("Buildings").Create(&asset).Error; err != nil {
			return err
		}
		return tx.Model(&asset).Association("Buildings").Replace(buildings)
	}); err != nil {
		c.JSON(500, gin.H{"error": "Failed to create asset"})
		return
	}

	c.JSON(201, asset)
}

func (h *AssetHandler) UpdateAsset(c *gin.Context) {
	assetID := c.Param("asset_id")

	var asset Asset
	if err := h.DB.First(&asset, "id = ?", assetID).Error; err != nil {
		c.JSON(404, gin.H{"error": "Asset not found"})
		return
	}

	if err := c.ShouldBindJSON(&asset); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	var buildings []Building
	if asset.BuildingIDs != nil {
		var err error
		buildings, err = h.getBuildings(asset.BuildingIDs)
		if err != nil {
			c.JSON(400, gin.H{"error": "One or more buildings were not found"})
			return
		}
	}

	if err := h.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Omit("Buildings").Save(&asset).Error; err != nil {
			return err
		}
		if asset.BuildingIDs != nil {
			return tx.Model(&asset).Association("Buildings").Replace(buildings)
		}
		return nil
	}); err != nil {
		c.JSON(500, gin.H{"error": "Failed to update asset"})
		return
	}

	c.JSON(200, asset)
}

func (h *AssetHandler) DeleteAsset(c *gin.Context) {
	assetID := c.Param("asset_id")

	if err := h.DB.Delete(&Asset{}, "id = ?", assetID).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete asset"})
		return
	}

	c.JSON(200, gin.H{"message": "Asset deleted successfully"})
}

func (h *AssetHandler) getBuildings(ids []uint) ([]Building, error) {
	if len(ids) == 0 {
		return []Building{}, nil
	}

	var buildings []Building
	if err := h.DB.Where("id IN ?", ids).Find(&buildings).Error; err != nil {
		return nil, err
	}
	if len(buildings) != len(ids) {
		return nil, gorm.ErrRecordNotFound
	}

	return buildings, nil
}
