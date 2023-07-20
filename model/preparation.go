package model

import (
	"time"

	"gorm.io/gorm"
)

// SHPreparation 高中共同備課記錄表
type SHPreparation struct {
	BaseModel
	DurationBaseModel
	Subject    string    `gorm:"not null;" form:"subject" i18n:"學科"`
	Grade      string    `gorm:"not null;" form:"grade" i18n:"教學年級"`
	Name       string    `gorm:"not null;" form:"name" i18n:"教學單元"`
	Material   string    `gorm:"not null;" form:"material" i18n:"教材來源"`
	Course     Course    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" i18n:"課程"`
	CourseID   SQLBasePK `gorm:"not null; index"`
	Objectives string    `gorm:"not null;" form:"objectives" i18n:"教學目標"`
	Experience string    `gorm:"not null;" form:"experience" i18n:"學生經驗"`
	Activity   string    `gorm:"not null;" form:"activity" i18n:"教學活動"`
	Evaluation string    `gorm:"not null;" form:"evaluation" i18n:"學生學習成效評估方式"`
}

// Save 儲存高中共同備課記錄表
func (object *SHPreparation) Save() error {
	return db.Save(object).Error
}

// GetTimeString get 高中共同備課記錄表的時間字串
func (object *SHPreparation) GetTimeString() (string, string, string) {
	return object.Date, object.StartTime, object.EndTime
}

// SetStartTime set 高中共同備課記錄表的開始時間
func (object *SHPreparation) SetStartTime(t time.Time) {
	object.Start = t
}

// SetEndTime set 高中共同備課記錄表的結束時間
func (object *SHPreparation) SetEndTime(t time.Time) {
	object.End = t
}

// GetSHPreparationByCourseID 依照課程編號取得共同備課記錄表
func GetSHPreparationByCourseID(courseID *SQLBasePK, object *SHPreparation) error {
	object.CourseID = *courseID
	return db.Joins("Course").Preload("Course.User").
		Where(object).First(object).Error
}

// AfterFind is a hook to format the date and time
func (object *SHPreparation) AfterFind(tx *gorm.DB) (err error) {
	object.Date = object.Start.Format(dateFormat)
	object.StartTime = object.Start.Format(pureTimeFormat)
	object.EndTime = object.End.Format(pureTimeFormat)
	return
}
