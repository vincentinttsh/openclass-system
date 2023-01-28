package model

import (
	"time"
	"vincentinttsh/openclass-system/pkg/tool"

	"gorm.io/gorm"
)

// Course is the base class model
type Course struct {
	BaseModel
	Name           string    `gorm:"not null;"`
	Classroom      string    `gorm:"not null;"`
	Department     string    `gorm:"not null;check:department in ('sh','jh')"`
	Start          time.Time `gorm:"index; type:date"`
	End            time.Time `gorm:"index; type:date check (\"end\" > \"start\")"`
	AttendPassword string    `gorm:"not null; type:char(6)"`
	UserID         SQLBasePK `gorm:"not null; index"`
	User           User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

// BeforeCreate set attend password
func (object *Course) BeforeCreate(tx *gorm.DB) (err error) {
	object.AttendPassword = tool.RandString(6)
	return
}

// Create create a class
func (object *Course) Create() error {
	result := db.Create(object)

	return result.Error
}

// Update update a class
func (object *Course) Update() error {
	result := db.Save(&object)

	return result.Error
}

// GetAllCourses get all class
func GetAllCourses(courses *[]Course) error {
	result := db.Model(Course{}).Joins("User").
		Where("start > ?", time.Now()).Order("start").Find(courses)

	return result.Error
}

// GetUserCourses get all class of single user
func GetUserCourses(userID *SQLBasePK, courses *[]Course) error {
	result := db.Model(Course{}).Joins("User").
		Where(&Course{UserID: *userID}).Order("`courses`.`id` DESC").Find(courses)

	return result.Error
}

// GetCourse get a class
func GetCourse(classID *SQLBasePK, class *Course, prefetch bool) error {
	var result *gorm.DB
	if prefetch {
		result = db.Model(Course{}).Joins("User").First(class, classID)
	} else {
		result = db.First(class, classID)
	}

	return result.Error
}
