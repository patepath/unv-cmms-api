package main

import (
	"fmt"
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
	ID             uint            `gorm:"primaryKey" json:"id"`
	Code           string          `gorm:"uniqueIndex;not null;size:100" json:"code"`
	Category       Category        `gorm:"foreignKey:CategoryID" json:"category"`
	CategoryID     uint            `gorm:"not null" json:"category_id"`
	Model          string          `gorm:"size:100" json:"model"`
	BrandName      string          `gorm:"size:100" json:"brand_name"`
	Location       Location        `gorm:"foreignKey:LocationID" json:"location"`
	LocationID     uint            `gorm:"not null" json:"location_id"`
	Area           string          `gorm:"size:255" json:"area"`
	Floor          string          `gorm:"size:255" json:"floor"`
	Owner          string          `gorm:"size:100" json:"owner"`
	Condition      uint            `gorm:"not null" json:"condition"`
	PurchaseDate   time.Time       `json:"purchase_date"`
	WarrantyEnd    time.Time       `json:"warranty_end"`
	LastService    time.Time       `json:"last_service"`
	NextService    time.Time       `json:"next_service"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
	DeletedAt      gorm.DeletedAt  `gorm:"index" json:"-"`
	ServiceRecords []ServiceRecord `gorm:"foreignKey:AssetID" json:"service_records,omitempty"`
}

type Category struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"size:100;not null" json:"name"`
}

type Location struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"size:255;not null" json:"name"`
}

type ServiceRecord struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	AssetID     uint           `gorm:"not null" json:"assetId"`
	ServiceDate time.Time      `json:"serviceDate,omitempty"`
	Details     string         `gorm:"type:text" json:"details"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func main() {
	var version = "api/v1"
	fmt.Println("Starting...")

	err := godotenv.Load(".env")
	if err != nil {
		log.Panic(err)
	}

	dsn := os.Getenv("dsn")

	var db *gorm.DB

	if db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{}); err != nil {
		log.Panic(err)

	} else {
		db.AutoMigrate(&Asset{}, &ServiceRecord{}, &Category{}, &Location{})
	}

	router := gin.Default()
	router.Use(CORSMiddleware())
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.SetTrustedProxies(nil)

	router.GET("/health", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	assetHandler := &AssetHandler{DB: db}
	assetGroup := router.Group("/" + version + "/asset")
	{
		assetGroup.GET("/", assetHandler.GetAssets)
		assetGroup.GET("/:id", assetHandler.GetAsset)
		assetGroup.POST("/", assetHandler.CreateAsset)
		assetGroup.PUT("/:id", assetHandler.UpdateAsset)
		//		assetGroup.DELETE("/:id", assetHandler.DeleteAsset)
	}

	serviceRecordHandler := &ServiceRecordHandler{DB: db}
	serviceRecordGroup := router.Group("/" + version + "/serviceRecord")
	{
		serviceRecordGroup.GET("/", serviceRecordHandler.GetServiceRecords)
		serviceRecordGroup.GET("/:id", serviceRecordHandler.GetServiceRecord)
		serviceRecordGroup.POST("/", serviceRecordHandler.CreateServiceRecord)
		serviceRecordGroup.PUT("/:id", serviceRecordHandler.UpdateServiceRecord)
		//		serviceRecordGroup.DELETE("/:id", serviceRecordHandler.DeleteServiceRecord)
	}

	locationHandler := &LocationHandler{DB: db}
	locationGroup := router.Group("/" + version + "/location")
	{
		locationGroup.GET("/", locationHandler.GetLocations)
		locationGroup.GET("/:id", locationHandler.GetLocationByID)
		locationGroup.POST("/", locationHandler.CreateLocation)
		locationGroup.PUT("/:id", locationHandler.UpdateLocation)
		// locationGroup.DELETE("/:id", locationHandler.DeleteLocation)
	}

	categoryHandler := &CategoryHandler{DB: db}
	categoryGroup := router.Group("/" + version + "/category")
	{
		categoryGroup.GET("/", categoryHandler.GetCategories)
		categoryGroup.GET("/:id", categoryHandler.GetCategoryByID)
		categoryGroup.POST("/", categoryHandler.CreateCategory)
		categoryGroup.PUT("/:id", categoryHandler.UpdateCategory)
		// categoryGroup.DELETE("/:id", categoryHandler.DeleteCategory)
	}

	if os.Getenv("production") == "true" {
		router.RunTLS(os.Getenv("api_host"), "/etc/nginx/ssl/cloudasprompt_com.crt", "/etc/nginx/ssl/cloudasprompt_com.key") // listen and serve on
	} else {
		router.Run(os.Getenv("api_host"))
	}
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
