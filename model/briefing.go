package model

import (
	"time"

	"gorm.io/gorm"
)

// SHBriefing 高中共同議課設計表
type SHBriefing struct {
	BaseModel
	Course          Course    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" i18n:"課程"`
	CourseID        SQLBasePK `gorm:"not null; index"`
	Observer        string    `gorm:"not null;" form:"observer" i18n:"觀察者"`
	Start           time.Time `gorm:"not null;"`
	Affirmation     string    `gorm:"not null;" form:"affirmation" i18n:"肯定教學表現"`
	GuideDiscussion string    `gorm:"not null;" form:"guideDiscussion" i18n:"引導討論教學表現"`
	Judgment        string    `gorm:"not null;" form:"judgment" i18n:"判斷表現程度"`
	Suggestion      string    `gorm:"not null;" form:"suggestion" i18n:"改進學生學習能力提供之建議"`
	Growth          string    `gorm:"not null;" form:"growth" i18n:"協助擬定成長活動"`
	Date            string    `gorm:"-:all" form:"date" i18n:"議課日期"`
	Time            string    `gorm:"-:all" form:"time" i18n:"議課時間"`
}

// Save 儲存高中共同議課設計表
func (object *SHBriefing) Save() error {
	return db.Save(object).Error
}

// AfterFind is a hook to format the date and time
func (object *SHBriefing) AfterFind(tx *gorm.DB) (err error) {
	object.Date = object.Start.Format(dateFormat)
	object.Time = object.Start.Format(pureTimeFormat)
	return
}

// GetSHBriefingByCourseID 依照課程編號取得共同備課記錄表
func GetSHBriefingByCourseID(courseID *SQLBasePK, object *SHBriefing) error {
	object.CourseID = *courseID
	return db.Joins("Course").Preload("Course.User").
		Where(object).First(object).Error
}
