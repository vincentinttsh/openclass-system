package model

import (
	"time"
	"vincentinttsh/openclass-system/pkg/tool"

	"gorm.io/gorm"
)

// BaseClass is the base class model
type BaseClass struct {
	BaseModel
	Name           string    `gorm:"not null;"`
	Classroom      string    `gorm:"not null;"`
	Start          time.Time `gorm:"index; type:date"`
	End            time.Time `gorm:"index; type:date check (\"end\" > \"start\")"`
	AttendPassword string    `gorm:"not null; type:char(6)"`
	TeacherID      uint      `gorm:"not null; index"`
	Teacher        User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;foreignKey:TeacherID"`
}

// BeforeCreate set attend password
func (object *BaseClass) BeforeCreate(tx *gorm.DB) (err error) {
	object.AttendPassword = tool.RandString(6)
	return
}

// Create create a class
func (object *BaseClass) Create() error {
	result := db.Create(&object)

	return result.Error
}