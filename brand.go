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
	if err := h.DB.Find(&brands).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, brands)
}
