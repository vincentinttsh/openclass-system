package view

const (
	dateFormat     string = "02/01/2006"
	pureTimeFormat string = "03:04 PM"
	timeFormat     string = "02/01/200603:04 PM"
)

var departmentChoice = map[string]string{
	"sh": "高中部",
	"jh": "國中部",
}
var subjectChoice = map[string]string{
	"chinese": "國文",
	"english": "英文",
	"math":    "數學",
	"science": "自然",
	"social":  "社會",
	"other":   "其他",
}

type baseClassStruct struct {
	Name      string `form:"className" validate:"required" i18n:"課程名稱"`
	Classroom string `form:"classroom" validate:"required" i18n:"上課教室"`
	Date      string `form:"date" validate:"required" i18n:"上課日期"`
	StartTime string `form:"startTime" validate:"required" i18n:"開始時間"`
	EndTime   string `form:"endTime" validate:"required" i18n:"結束時間"`
}

type registerUser struct {
	Username   string `form:"username" i18n:"姓名"`
	Email      string `form:"email" i18n:"Email address"`
	Department string `form:"department" i18n:"部門" validate:"required,oneof=sh jh"`
	Subject    string `form:"subject" i18n:"授課科目" validate:"required,oneof=chinese english math science social other"`
}
