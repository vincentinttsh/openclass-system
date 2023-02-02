package view

import (
	"errors"
	"strconv"
	"time"
	"vincentinttsh/openclass-system/model"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// CreateOpenClass create openclass class(POST, GET)
func CreateOpenClass(c *fiber.Ctx) error {
	var userID model.SQLBasePK = model.SQLBasePK(c.Locals("id").(float64))
	var form model.Course
	var valid bool
	var err error
	var template string = "class/form"
	var bind fiber.Map = c.Locals("bind").(fiber.Map)

	bind["method"] = "新增"
	bind["permissions"] = "edit"
	bind["csrf_token"] = c.Locals("csrf_token")
	bind["form"] = &form

	if c.Method() == fiber.MethodGet {
		form.Date = time.Now().Format(dateFormat)
		form.StartTime = time.Now().Format(pureTimeFormat)
		form.EndTime = time.Now().Add(time.Hour).Format(pureTimeFormat)
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
	form.Start, _ = time.Parse(timeFormat, form.Date+form.StartTime)
	form.End, _ = time.Parse(timeFormat, form.Date+form.EndTime)

	// 轉換時區
	form.Start = form.Start.Add(timeOffset).Local()
	form.End = form.End.Add(timeOffset).Local()

	if form.End.Before(form.Start) {
		bind["messages"] = []msgStruct{
			createMsg(warnMsgLevel, "結束時間不得早於開始時間"),
		}
		return c.Status(fiber.StatusBadRequest).Render(template, bind)
	}

	form.Department = c.Locals("department").(string)
	form.UserID = userID

	if err = form.Save(); err != nil {
		return dbWriteError(c, err, template, &bind)
	}

	return c.RedirectToRoute("home", fiber.Map{
		"queries": map[string]string{
			"status": "create_success",
		},
	})
}

// ListUserOpenClass list user openclass class(GET)
func ListUserOpenClass(c *fiber.Ctx) error {
	var userID = model.SQLBasePK(c.Locals("id").(float64))
	var courses []model.Course
	var err error
	var template string = "class/list"
	var bind fiber.Map = c.Locals("bind").(fiber.Map)
	bind["baseURL"] = baseURL
	bind["courses"] = &courses

	err = model.GetUserCourses(&userID, &courses)
	if err != nil {
		return dbReadError(c, err, template, &bind)
	}

	var subCount int64
	var attendeesCount int64
	for i, v := range courses {
		attendeesCount, subCount, err = v.GetAttendeesCount()
		if err != nil {
			return dbReadError(c, err, template, &bind)
		}
		courses[i].AttendCount = attendeesCount
		courses[i].SubCount = subCount
	}

	statusBinding(c, &bind)

	return c.Status(fiber.StatusOK).Render(template, bind)
}

// GetOfModifyOpenClass openclass class(GET, POST)
func GetOfModifyOpenClass(c *fiber.Ctx) error {
	var classID model.SQLBasePK = 0
	var id uint64
	var form model.Course
	var valid bool
	var err error
	var template string = "class/form"
	var bind fiber.Map = c.Locals("bind").(fiber.Map)

	if id, err = strconv.ParseUint(c.Params("id"), 10, 64); err != nil {
		return notFound(c)
	}
	classID = model.SQLBasePK(id)

	bind["method"] = "修改"
	bind["permissions"] = "edit"
	bind["csrf_token"] = c.Locals("csrf_token")
	bind["form"] = &form

	err = model.GetCourse(&classID, &form, false)
	if err == gorm.ErrRecordNotFound {
		return notFound(c)
	}
	if err != nil {
		return dbReadError(c, err, template, &bind)
	}

	if c.Method() == "GET" {
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
	form.Start, _ = time.Parse(timeFormat, form.Date+form.StartTime)
	form.End, _ = time.Parse(timeFormat, form.Date+form.EndTime)
	form.Start = form.Start.Add(timeOffset).Local()
	form.End = form.End.Add(timeOffset).Local()

	if form.End.Before(form.Start) {
		bind["messages"] = []msgStruct{
			createMsg(warnMsgLevel, "結束時間不得早於開始時間"),
		}
		return c.Status(fiber.StatusBadRequest).Render(template, bind)
	}

	if err = form.Save(); err != nil {
		bind["messages"] = []msgStruct{
			createMsg(errMsgLevel, serverErrorMsg),
		}
		sugar.Errorln(err)
		return c.Status(fiber.StatusInternalServerError).Render(template, bind)
	}

	return c.RedirectToRoute("my_class", fiber.Map{
		"queries": map[string]string{
			"status": "update_success",
		},
	})
}

// AttendClass 報到課程
func AttendClass(c *fiber.Ctx) error {
	var id uint64
	var classID model.SQLBasePK
	var userID model.SQLBasePK = model.SQLBasePK(c.Locals("id").(float64))
	var course model.Course
	var err error
	var msg string
	var template string = "error/msg"
	var bind fiber.Map = c.Locals("bind").(fiber.Map)
	var now = time.Now().Local()
	var passwd = c.Query("passwd", "")

	if id, err = strconv.ParseUint(c.Params("id"), 10, 64); err != nil {
		return notFound(c)
	}
	classID = model.SQLBasePK(id)

	err = model.GetCourse(&classID, &course, false)
	if err == gorm.ErrRecordNotFound {
		return notFound(c)
	}
	if err != nil {
		return dbReadError(c, err, template, &bind)
	}

	// 檢查密碼
	if passwd != course.AttendPassword {
		err = errors.New("密碼錯誤")
		msg = "請確認密碼是否正確"
	}
	// 檢查時間
	course.Start = course.Start.Add(allowTimeOffset)
	if now.Before(course.Start) {
		err = errors.New("課程尚未開始")
		msg = "請於 " + course.Start.Format("2006年01月02日 03:04 PM") + " 後進行報到"
	}
	if now.After(course.End) {
		err = errors.New("課程已結束")
		msg = "請於 " + course.End.Format("2006年01月02日 03:04 PM") + " 前進行報到"
	}
	if err != nil {
		bind["messages"] = []msgStruct{createMsgWithDetail(warnMsgLevel, err.Error(), msg)}
		return c.Status(fiber.StatusBadRequest).Render(template, bind)
	}

	err = course.Attend(&userID)
	if err != nil {
		if err == model.ErrAttendYourClass {
			bind["messages"] = []msgStruct{createMsg(warnMsgLevel, err.Error())}
			return c.Status(fiber.StatusBadRequest).Render(template, bind)
		}
		return serverError(c, err, template, &bind)
	}

	return c.RedirectToRoute("my_observation", fiber.Map{
		"queries": map[string]string{
			"status": "attend_success",
		},
	})
}

// ReserveClass 預約課程
func ReserveClass(c *fiber.Ctx) error {
	var id uint64
	var classID model.SQLBasePK
	var userID model.SQLBasePK = model.SQLBasePK(c.Locals("id").(float64))
	var course model.Course
	var err error
	var template string = "error/msg"
	var bind fiber.Map = c.Locals("bind").(fiber.Map)

	if id, err = strconv.ParseUint(c.Params("id"), 10, 64); err != nil {
		return notFound(c)
	}
	classID = model.SQLBasePK(id)

	err = model.GetCourse(&classID, &course, false)
	if err == gorm.ErrRecordNotFound {
		return notFound(c)
	}
	if err != nil {
		return dbReadError(c, err, template, &bind)
	}

	err = course.Reserve(&userID)
	if err != nil {
		if err == model.ErrAttendYourClass {
			bind["messages"] = []msgStruct{createMsg(warnMsgLevel, err.Error())}
			return c.Status(fiber.StatusBadRequest).Render(template, bind)
		}
		return serverError(c, err, template, &bind)
	}

	return c.RedirectToRoute("my_observation", fiber.Map{
		"queries": map[string]string{
			"status": "reserve_success",
		},
	})
}
