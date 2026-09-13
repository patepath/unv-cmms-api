package main

import (
	"gorm.io/gorm"
)

type AssetHandler struct {
	DB *gorm.DB
}
