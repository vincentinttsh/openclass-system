package model

import (
	"errors"
	"fmt"
	"time"
	"vincentinttsh/openclass-system/internal/tool"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ErrAttendYourClass 不能報到/預約自己的課程
var ErrAttendYourClass = errors.New("您不能報到/預約自己的課程")
var departmentChoice = map[string]string{
	"sh": "高中部",
	"jh": "國中部",
}
var calendar string = `
http://www.google.com/calendar/event?action=TEMPLATE&text=%s公開授課（%s)&dates=%s/%s&details=課程名稱：%s
`

// Course is the base class model
type Course struct {
	BaseModel
	Name           string    `gorm:"not null;" form:"name" validate:"required" i18n:"課程名稱"`
	Classroom      string    `gorm:"not null;" form:"classroom" validate:"required" i18n:"上課教室"`
	Department     string    `gorm:"not null;check:department in ('sh','jh')"`
	Start          time.Time `gorm:"index; type:date"`
	End            time.Time `gorm:"index; type:date check (\"end\" > \"start\")"`
	AttendPassword string    `gorm:"not null; type:char(6)"`
	UserID         SQLBasePK `gorm:"not null; index"`
	User           User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Date           string    `gorm:"-:all" form:"date" validate:"required" i18n:"上課日期"`
	StartTime      string    `gorm:"-:all" form:"startTime" validate:"required" i18n:"開始時間"`
	EndTime        string    `gorm:"-:all" form:"endTime" validate:"required" i18n:"結束時間"`
	DateCH         string    `gorm:"-:all"`
	Duration       string    `gorm:"-:all"`
	Calendar       string    `gorm:"-:all"`
	AttendCount    int64     `gorm:"-:all"`
	SubCount       int64     `gorm:"-:all"`
}

// CourseReservation 課程預約 model
type CourseReservation struct {
	CreatedAt time.Time
	CourseID  SQLBasePK `gorm:"not null;index;primaryKey"`
	Course    Course    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	UserID    SQLBasePK `gorm:"not null;index;primaryKey"`
	User      User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Attended  bool      `gorm:"not null;index;default:false"`
}

// BeforeCreate set attend password
func (object *Course) BeforeCreate(tx *gorm.DB) (err error) {
	object.AttendPassword = tool.RandString(6)
	return
}

// BeforeCreate check if the user is the course owner 不能報到自己的課程
func (object *CourseReservation) BeforeCreate(tx *gorm.DB) (err error) {
	var course Course
	err = GetCourse(&object.CourseID, &course, false)
	if err != nil {
		return
	}
	if course.UserID == object.UserID {
		err = ErrAttendYourClass
	}
	return
}

// AfterFind is a hook to format the date and time
func (object *Course) AfterFind(tx *gorm.DB) (err error) {
	object.DateCH = object.Start.Format("2006年01月02日")
	object.Date = object.Start.Format(dateFormat)
	object.StartTime = object.Start.Format(pureTimeFormat)
	object.EndTime = object.End.Format(pureTimeFormat)
	// User is prefetch
	if object.User.ID != 0 {
		object.Duration = object.Start.Format("15:04") + " ~ " + object.End.Format("15:04")
		object.Calendar = fmt.Sprintf(calendar,
			departmentChoice[object.User.Department],
			object.User.Name,
			object.Start.Format("20060102T150405"),
			object.End.Format("20060102T150405"),
			object.Name,
		) + "%0A" + fmt.Sprintf("授課老師：%s&location=%s&trp=false", object.User.Name, object.Classroom)
	}
	return
}

// AfterFind is a hook to format the date and time
func (object *CourseReservation) AfterFind(tx *gorm.DB) (err error) {
	object.Course.AfterFind(tx)
	return
}

// Save save a class
func (object *Course) Save() error {
	return db.Save(&object).Error
}

// Attend 老師報到課程
func (object *Course) Attend(userID *SQLBasePK) error {
	return db.Clauses(clause.OnConflict{
		DoUpdates: clause.Assignments(map[string]interface{}{"attended": true}),
	}).Create(&CourseReservation{
		CourseID: object.ID,
		UserID:   *userID,
		Attended: true,
	}).Error
}

// Reserve 老師預約課程
func (object *Course) Reserve(userID *SQLBasePK) (err error) {
	return db.Clauses(clause.OnConflict{
		DoNothing: true,
	}).Create(&CourseReservation{
		CourseID: object.ID,
		UserID:   *userID,
		Attended: false,
	}).Error
}

// GetAttendees get the course attendees
func (object *Course) GetAttendees() (reservations []CourseReservation, err error) {
	err = db.Model(&CourseReservation{}).Where(&CourseReservation{
		CourseID: object.ID,
	}).Find(&reservations).Error
	return
}

// GetAttendeesCount get the number of course attendees
func (object *Course) GetAttendeesCount() (attendee int64, subscriber int64, err error) {
	err = db.Model(&CourseReservation{}).Where(&CourseReservation{
		CourseID: object.ID,
		Attended: true,
	}).Count(&attendee).Error
	if err != nil {
		return
	}
	err = db.Model(&CourseReservation{}).Where(&CourseReservation{
		CourseID: object.ID,
	}).Not(&CourseReservation{
		Attended: true,
	}).Count(&subscriber).Error

	return
}

// GetAllCourses get all class
func GetAllCourses(courses *[]Course) error {
	return db.Joins("User").
		Where("end > ?", time.Now()).Order("start").Find(courses).Error
}

// GetUserCourses get all class of single user
func GetUserCourses(userID *SQLBasePK, courses *[]Course) error {
	return db.Joins("User").Where(
		&Course{UserID: *userID}).Order("`courses`.`id` DESC").Find(courses).Error
}

// GetUserObserveCourses 取得使用者觀課或預約的課程
func GetUserObserveCourses(userID *SQLBasePK, reservations *[]CourseReservation) error {
	return db.Joins("Course").Where(
		&CourseReservation{UserID: *userID},
	).Order("`course_reservations`.`created_at` DESC").Find(reservations).Error
}

// GetCourseObserve get the course attendees
func GetCourseObserve(courseID *SQLBasePK, reservations *[]CourseReservation) error {
	return db.Joins("User").Where(&CourseReservation{
		CourseID: *courseID,
		Attended: true,
	}).Order("`course_reservations`.`created_at` DESC").Find(reservations).Error
}

// GetCourse get a class
func GetCourse(classID *SQLBasePK, class *Course, prefetch bool) error {
	class.ID = *classID
	if prefetch {
		return db.Joins("User").First(class).Error
	}
	return db.First(class).Error
}
