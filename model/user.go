package model

import "gorm.io/gorm"

// User user model
type User struct {
	CreatedAt      int64
	UpdatedAt      int64
	DeletedAt      gorm.DeletedAt `gorm:"index"`
	ID             string         `gorm:"not null;primary_key"`
	Email          string         `gorm:"not null;unique_index"`
	Password       string         `gorm:"not null"`
	Name           string         `gorm:"not null"`
	Admin          bool           `gorm:"not null;index;default:false"`
	SuperAdmin     bool           `gorm:"not null;index;default:false"`
	OrganizationID uint           `gorm:"not null;index"`
	Organization   Organization   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
}

// Organization organization model (usually a school)
type Organization struct {
	BaseModel
	Stage   string `gorm:"index;not null;size:255"`
	Name    string `gorm:"index;not null;size:255"`
	Address string `gorm:"not null;size:255"`
}
