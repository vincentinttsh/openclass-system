package view

import (
	"fmt"
	"strconv"
	"time"
	"vincentinttsh/openclass-system/model"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// CreateOpenClass create openclass class(POST, GET)
func CreateOpenClass(c *fiber.Ctx) error {
	var userID model.SQLBasePK = model.SQLBasePK(c.Locals("id").(float64))
	var template string = "class/form"
	var form baseClassStruct
	var course model.Course
	var valid bool
	var err error
	var startTime time.Time
	var endTime time.Time
	var bind fiber.Map = c.Locals("bind").(fiber.Map)

	form = baseClassStruct{
		Date:      time.Now().Format(dateFormat),
		StartTime: time.Now().Format(pureTimeFormat),
		EndTime:   time.Now().Add(time.Hour).Format(pureTimeFormat),
	}
	bind["method"] = "新增"
	bind["permissions"] = "edit"
	bind["csrf_token"] = c.Locals("csrf_token")
	bind["form"] = form

	if c.Method() == fiber.MethodGet {
		return c.Status(fiber.StatusOK).Render(template, bind)
	}

	valid, err = formValidate(c, &bind, &form, &template)
	if !valid {
		return err
	}

	value := []string{form.Date, form.Date + form.StartTime, form.Date + form.EndTime}
	parseKey := []string{"日期", "開始時間", "結束時間"}
	parseFormat := []string{dateFormat, timeFormat, timeFormat}
	for i := 0; i < len(value); i++ {
		if _, err = time.Parse(parseFormat[i], value[i]); err != nil {
			badRequest(c, parseKey[i]+formatErrorMsg, template, &bind)
		}
	}
	startTime, _ = time.Parse(timeFormat, form.Date+form.StartTime)
	endTime, _ = time.Parse(timeFormat, form.Date+form.EndTime)

	// 轉換時區
	startTime = startTime.Add(timeOffset)
	endTime = endTime.Add(timeOffset)

	if endTime.Before(startTime) {
		bind["messages"] = []msgStruct{
			createMsg(warnMsgLevel, "結束時間不得早於開始時間"),
		}
		return c.Status(fiber.StatusBadRequest).Render(template, bind)
	}

	course = model.Course{
		Name:       form.Name,
		Department: c.Locals("department").(string),
		Classroom:  form.Classroom,
		Start:      startTime.Local(),
		End:        endTime.Local(),
		UserID:     userID,
	}

	if err = course.Create(); err != nil {
		return dbWriteError(c, err, template, &bind)
	}

	return c.Redirect("/?status=create_success")
}

// ListUserOpenClass list user openclass class(GET)
func ListUserOpenClass(c *fiber.Ctx) error {
	var userID = model.SQLBasePK(c.Locals("id").(float64))
	var template string = "class/list"
	var data []model.Course
	var courses []courseBind
	var err error
	var bind fiber.Map = c.Locals("bind").(fiber.Map)
	bind["baseURL"] = baseURL

	err = model.GetUserCourses(&userID, &data)
	if err != nil {
		sugar.Errorw("Get all class error", "error", err)
		bind["messages"] = []msgStruct{
			createMsg(errMsgLevel, "取得課程資料時發生錯誤"),
		}
		return c.Status(fiber.StatusInternalServerError).Render(template, bind)
	}

	courses = make([]courseBind, len(data))
	var calendar string = "http://www.google.com/calendar/event?action=TEMPLATE&text=%s公開授課（"
	calendar += "%s)&dates=%s/%s&details=課程名稱：%s"
	for i, v := range data {
		courses[i] = courseBind{
			ClassID:   v.ID,
			ClassName: v.Name,
			Date:      v.Start.Format("2006年01月02日"),
			Passwd:    v.AttendPassword,
			Calendar: fmt.Sprintf(calendar,
				departmentChoice[v.User.Department],
				v.User.Name,
				v.Start.Format("20060102T150405"),
				v.End.Format("20060102T150405"),
				v.Name,
			) + "%0A" + fmt.Sprintf("授課老師：%s&location=%s&trp=false", v.User.Name, v.Classroom),
		}
	}

	bind["courses"] = &courses

	return c.Status(fiber.StatusOK).Render(template, bind)
}

// GetOfModifyOpenClass openclass class(GET, POST)
func GetOfModifyOpenClass(c *fiber.Ctx) error {
	var classID model.SQLBasePK = 0
	var id uint64
	var template string = "class/form"
	var data model.Course
	var bind fiber.Map
	var form baseClassStruct
	var valid bool
	var startTime time.Time
	var endTime time.Time
	var err error

	if id, err = strconv.ParseUint(c.Params("id"), 10, 64); err != nil {
		return notFound(c)
	}
	classID = model.SQLBasePK(id)
	bind = c.Locals("bind").(fiber.Map)
	bind["method"] = "修改"
	bind["permissions"] = "edit"
	bind["csrf_token"] = c.Locals("csrf_token")

	err = model.GetCourse(&classID, &data, false)
	if err == gorm.ErrRecordNotFound {
		return notFound(c)
	}
	if err != nil {
		return dbReadError(c, err, template, &bind)
	}

	if c.Method() == "GET" {
		form = baseClassStruct{
			Name:      data.Name,
			Classroom: data.Classroom,
			Date:      data.Start.Format(dateFormat),
			StartTime: data.Start.Format(pureTimeFormat),
			EndTime:   data.End.Format(pureTimeFormat),
		}
		bind["form"] = form

		return c.Status(fiber.StatusOK).Render(template, bind)
	}

	valid, err = formValidate(c, &bind, &form, &template)
	if !valid {
		return err
	}

	value := []string{form.Date, form.Date + form.StartTime, form.Date + form.EndTime}
	parseKey := []string{"日期", "開始時間", "結束時間"}
	parseFormat := []string{dateFormat, timeFormat, timeFormat}
	for i := 0; i < len(value); i++ {
		if _, err = time.Parse(parseFormat[i], value[i]); err != nil {
			bind["messages"] = []msgStruct{
				createMsg(warnMsgLevel, parseKey[i]+formatErrorMsg),
			}
			return c.Status(fiber.StatusBadRequest).Render(template, bind)
		}
	}
	startTime, _ = time.Parse(timeFormat, form.Date+form.StartTime)
	endTime, _ = time.Parse(timeFormat, form.Date+form.EndTime)

	// 轉換時區
	startTime = startTime.Add(timeOffset)
	endTime = endTime.Add(timeOffset)

	if endTime.Before(startTime) {
		bind["messages"] = []msgStruct{
			createMsg(warnMsgLevel, "結束時間不得早於開始時間"),
		}
		return c.Status(fiber.StatusBadRequest).Render(template, bind)
	}

	data.Name = form.Name
	data.Classroom = form.Classroom
	data.Start = startTime.Local()
	data.End = endTime.Local()

	if err = data.Update(); err != nil {
		bind["messages"] = []msgStruct{
			createMsg(errMsgLevel, serverErrorMsg),
		}
		sugar.Errorln(err)
		return c.Status(fiber.StatusInternalServerError).Render(template, bind)
	}

	return c.Redirect("/?status=update_success")
}
