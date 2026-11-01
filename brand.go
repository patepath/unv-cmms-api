package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BrandHandler struct {
	DB *gorm.DB
}

func (h *BrandHandler) GetBrands(c *gin.Context) {
	var brands []Brand
	if err := h.DB.Order("name").Find(&brands).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, brands)
}

func (h *BrandHandler) GetBrandByID(c *gin.Context) {
	id := c.Param("id")
	var brand Brand
	if err := h.DB.First(&brand, id).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, brand)
}

func (h *BrandHandler) CreateBrand(c *gin.Context) {
	var brand Brand
	if err := c.ShouldBindJSON(&brand); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := h.DB.Create(&brand).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, brand)
}

func (h *BrandHandler) UpdateBrand(c *gin.Context) {
	id := c.Param("id")
	var brand Brand
	if err := h.DB.First(&brand, id).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if err := c.ShouldBindJSON(&brand); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := h.DB.Save(&brand).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, brand)
}

func (h *BrandHandler) DeleteBrand(c *gin.Context) {
	id := c.Param("id")
	if err := h.DB.Delete(&Brand{}, id).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(204, nil)
}
