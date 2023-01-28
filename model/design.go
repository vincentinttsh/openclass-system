package model

import (
	"time"

	"gorm.io/gorm"
)

// SHDesign is the model for 高中課程教學活動設計表
type SHDesign struct {
	BaseModel
	Name       string           `gorm:"not null;" form:"className" i18n:"單元名稱"`
	Material   string           `gorm:"not null;" form:"material" i18n:"教材來源"`
	Class      string           `gorm:"not null;" form:"class" i18n:"授課班級"`
	Start      time.Time        `gorm:"index; type:date"`
	End        time.Time        `gorm:"index; type:date check (\"end\" > \"start\")"`
	Objectives string           `gorm:"not null;" form:"objectives" i18n:"教學目標"`
	Background string           `gorm:"not null;" form:"background" i18n:"學生學習背景分析"`
	Method     string           `gorm:"not null;" form:"method" i18n:"教學方法"`
	Resource   string           `gorm:"not null;" form:"resource" i18n:"教學資源"`
	Reference  string           `gorm:"not null;" form:"reference" i18n:"參考資料"`
	CourseID   SQLBasePK        `gorm:"not null; index"`
	Course     Course           `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Date       string           `gorm:"-:all" form:"date" i18n:"教學日期"`
	StartTime  string           `gorm:"-:all" form:"startTime" i18n:"教學開始時間"`
	EndTime    string           `gorm:"-:all" form:"endTime" i18n:"教學結束時間"`
	Details    []SHDesignDetail `gorm:"-:all" form:"detail"`
}

// SHDesignDetail is the model for 高中課程教學活動設計表的詳細資料
type SHDesignDetail struct {
	BaseModel
	Refreshed       bool      `gorm:"-:all"`
	SHDesignID      SQLBasePK `gorm:"not null; index"`
	SHDesign        SHDesign  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:SHDesignID"`
	Objectives      string    `gorm:"not null;" form:"objectives" i18n:"教學目標"`
	TeacherActivity string    `gorm:"not null;" form:"teacherActivity" i18n:"教師活動"`
	StudentActivity string    `gorm:"not null;" form:"studentActivity" i18n:"學生活動"`
	Aid             string    `gorm:"not null;" form:"aid" i18n:"使用教具"`
	Score           string    `gorm:"not null;" form:"score" i18n:"評量方式"`
	Method          string    `gorm:"not null;" form:"method" i18n:"教學方法"`
	Allocation      string    `gorm:"not null;" form:"allocation" i18n:"時間分配"`
}

// Create create a 高中課程教學活動設計表的詳細資料
func (object *SHDesignDetail) Create() error {
	result := db.Create(&object)

	return result.Error
}

// Create create a 高中課程教學活動設計表
func (object *SHDesign) Create() error {
	result := db.Create(object)

	return result.Error
}

// Update update a 高中課程教學活動設計表
func (object *SHDesign) Update() error {
	result := db.Save(object)

	return result.Error
}

// Update update a 高中課程教學活動設計表的詳細資料
func (object *SHDesignDetail) Update() error {
	var result *gorm.DB
	if object.ID == 0 {
		result = db.Save(object)
	} else {
		result = db.Model(object).Updates(object)
	}

	return result.Error
}

// Delete delete a 高中課程教學活動設計表詳細資料
func (object *SHDesignDetail) Delete() error {
	result := db.Delete(object)

	return result.Error
}

// GetSHDesignByCourseID get a 高中課程教學活動設計表 by courseID
func GetSHDesignByCourseID(courseID *SQLBasePK, object *SHDesign) error {
	result := db.Model(SHDesign{}).Joins("Course").Preload("Course.User").Where(&SHDesign{
		CourseID: *courseID,
	}).First(&object)

	return result.Error
}

// GetCourseDesignDetails 取得高中課程教學活動設計表的詳細資料
func GetCourseDesignDetails(SHDesignID *SQLBasePK, objects *[]SHDesignDetail) error {
	result := db.Where(&SHDesignDetail{
		SHDesignID: *SHDesignID,
	}).Find(&objects).Order("id")

	return result.Error
}

// AfterFind is a hook to format the date and time
func (object *SHDesign) AfterFind(tx *gorm.DB) (err error) {
	object.Date = object.Start.Format(dateFormat)
	object.StartTime = object.Start.Format(pureTimeFormat)
	object.EndTime = object.End.Format(pureTimeFormat)
	err = GetCourseDesignDetails(&object.ID, &object.Details)
	if len(object.Details) == 0 {
		object.Details = make([]SHDesignDetail, 1)
	}
	return
}
