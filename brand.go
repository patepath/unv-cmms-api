package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BrandHandler struct {
	DB *gorm.DB
}

func (h *BrandHandler) GetBrands(c *gin.Context) {
	var brands []Brand
	if err := h.DB.Preload("Category").Order("name").Find(&brands).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, brands)
}

func (h *BrandHandler) GetBrandsByCategory(c *gin.Context) {
	categoryID := c.Param("category_id")
	var brands []Brand
	if err := h.DB.Preload("Category").
		Where("category_id = ?", categoryID).
		Order("name").
		Find(&brands).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, brands)
}

func (h *BrandHandler) GetBrandByID(c *gin.Context) {
	id := c.Param("id")
	var brand Brand
	if err := h.DB.Preload("Category").First(&brand, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Brand not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, brand)
}

func (h *BrandHandler) CreateBrand(c *gin.Context) {
	var brand Brand
	if err := c.ShouldBindJSON(&brand); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	brand.Category = Category{}
	if err := h.DB.Omit("Category").Create(&brand).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := h.DB.Preload("Category").First(&brand, brand.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, brand)
}

func (h *BrandHandler) UpdateBrand(c *gin.Context) {
	id := c.Param("id")
	var brand Brand
	if err := h.DB.First(&brand, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Brand not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var input struct {
		CategoryID *uint  `json:"category_id"`
		Name       string `json:"name"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	brand.CategoryID = input.CategoryID
	brand.Name = input.Name
	if err := h.DB.Model(&brand).Select("CategoryID", "Name").Updates(&brand).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	brand.Category = Category{}
	if err := h.DB.Preload("Category").First(&brand, brand.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, brand)
}

func (h *BrandHandler) DeleteBrand(c *gin.Context) {
	id := c.Param("id")
	if err := h.DB.Delete(&Brand{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
