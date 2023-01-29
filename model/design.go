package model

import (
	"time"

	"gorm.io/gorm"
)

// SHDesign is the model for 高中課程教學活動設計表
type SHDesign struct {
	BaseModel
	DurationBaseModel
	Name       string           `gorm:"not null;" form:"name" i18n:"單元名稱"`
	Material   string           `gorm:"not null;" form:"material" i18n:"教材來源"`
	Class      string           `gorm:"not null;" form:"class" i18n:"授課班級"`
	Objectives string           `gorm:"not null;" form:"objectives" i18n:"教學目標"`
	Background string           `gorm:"not null;" form:"background" i18n:"學生學習背景分析"`
	Method     string           `gorm:"not null;" form:"method" i18n:"教學方法"`
	Resource   string           `gorm:"not null;" form:"resource" i18n:"教學資源"`
	Reference  string           `gorm:"not null;" form:"reference" i18n:"參考資料"`
	CourseID   SQLBasePK        `gorm:"not null; index"`
	Course     Course           `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" i18n:"課程"`
	Details    []SHDesignDetail `gorm:"-:all" form:"detail"`
}

// SHDesignDetail is the model for 高中課程教學活動設計表的詳細資料
type SHDesignDetail struct {
	BaseModel
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

// Save 儲存高中課程教學活動設計表
func (object *SHDesign) Save() error {
	return db.Save(object).Error
}

// GetTimeString get 高中課程教學活動設計表的時間字串
func (object *SHDesign) GetTimeString() (string, string, string) {
	return object.Date, object.StartTime, object.EndTime
}

// SetStartTime set 高中課程教學活動設計表的開始時間
func (object *SHDesign) SetStartTime(t time.Time) {
	object.Start = t
}

// SetEndTime set 高中課程教學活動設計表的結束時間
func (object *SHDesign) SetEndTime(t time.Time) {
	object.End = t
}

// Save 儲存高中課程教學活動設計表的詳細資料
func (object *SHDesignDetail) Save() error {
	return db.Save(object).Error
}

// Delete delete a 高中課程教學活動設計表詳細資料
func (object *SHDesignDetail) Delete() error {
	return db.Delete(object).Error
}

// GetSHDesignByCourseID 依照課程編號取得高中課程教學活動設計表
func GetSHDesignByCourseID(courseID *SQLBasePK, object *SHDesign) error {
	object.CourseID = *courseID
	return db.Model(SHDesign{}).Joins("Course").Preload("Course.User").First(&object).Error
}

// GetCourseDesignDetails 取得高中課程教學活動設計表的詳細資料
func GetCourseDesignDetails(SHDesignID *SQLBasePK, objects *[]SHDesignDetail) error {
	return db.Where(&SHDesignDetail{
		SHDesignID: *SHDesignID,
	}).Find(&objects).Order("id").Error
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
