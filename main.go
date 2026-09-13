package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Asset struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	AssetCode    string         `gorm:"uniqueIndex;not null;size:100" json:"assetCode"`
	Name         string         `gorm:"not null;size:255" json:"name"`
	Description  string         `gorm:"type:text" json:"description"`
	SerialNumber string         `gorm:"size:100" json:"serialNumber"`
	Model        string         `gorm:"size:100" json:"model"`
	Manufacturer string         `gorm:"size:100" json:"manufacturer"`
	Location     string         `gorm:"size:255" json:"location"`
	Status       string         `gorm:"not null;size:50" json:"status"`
	PurchaseDate *time.Time     `json:"purchaseDate,omitempty"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Panic(err)
	}

	dsn := os.Getenv("dsn")

	var db *gorm.DB

	if db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{}); err != nil {
		log.Panic(err)

	} else {
		db.AutoMigrate(Asset{})
	}

	router := gin.Default()
	router.Use(CORSMiddleware())

	router.GET("/health", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router.Run(":" + port)
}

func CORSMiddleware() gin.HandlerFunc {
	allowOrigins := map[string]bool{
		"http://localhost:4200": true,
	}

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		if allowOrigins[origin] {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Cache-Control, X-Requested-With")
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
