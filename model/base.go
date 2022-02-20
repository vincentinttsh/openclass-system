package model

import (
	"vincentinttsh/openclass-system/config"

	"gorm.io/gorm"
)

var db *gorm.DB

// InitFunc is a function that initializes the model
func InitFunc(config *config.Config) {
	db = config.DB
	db.AutoMigrate(&Organization{})
	db.AutoMigrate(&User{})
}

// BaseModel base model
type BaseModel struct {
	CreatedAt int64
	UpdatedAt int64
	DeletedAt gorm.DeletedAt `gorm:"index"`
	ID        uint           `gorm:"not null;primaryKey"`
}
