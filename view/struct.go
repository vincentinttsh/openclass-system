package view

const (
	dateFormat string = "02/01/2006"
	timeFormat string = "02/01/200603:04 PM"
)

type baseClassStruct struct {
	Name      string `form:"className" validate:"required" i18n:"課程名稱"`
	Classroom string `form:"classroom" validate:"required" i18n:"上課教室"`
	Date      string `form:"date" validate:"required" i18n:"上課日期"`
	StartTime string `form:"startTime" validate:"required" i18n:"開始時間"`
	EndTime   string `form:"endTime" validate:"required" i18n:"結束時間"`
}
