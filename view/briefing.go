package view

import (
	"strconv"
	"time"
	"vincentinttsh/openclass-system/model"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func createCourseBriefingSH(c *fiber.Ctx) error {
	var classID model.SQLBasePK
	var id uint64
	var form model.SHBriefing
	var template string = "briefing/formSH"
	var bind fiber.Map = c.Locals("bind").(fiber.Map)
	var err error

	bind["method"] = "儲存"
	bind["permissions"] = "edit"
	bind["csrf_token"] = c.Locals("csrf_token")
	bind["form"] = &form

	if id, err = strconv.ParseUint(c.Params("id"), 10, 64); err != nil {
		return notFound(c)
	}
	classID = model.SQLBasePK(id)

	err = model.GetSHBriefingByCourseID(&classID, &form)
	switch err {
	case gorm.ErrRecordNotFound:
		var courseData model.Course
		var attendee []model.CourseReservation

		err = model.GetCourse(&classID, &courseData, true)
		switch err {
		case gorm.ErrRecordNotFound:
			return notFound(c)
		case nil:
			break
		default:
			return dbReadError(c, err, template, &bind)
		}
		err = model.GetCourseObserve(&classID, &attendee)
		if err != nil {
			return dbReadError(c, err, template, &bind)
		}
		form.Course = courseData
		form.Date = time.Now().Format(dateFormat)
		form.Time = time.Now().Format(pureTimeFormat)
		for i, v := range attendee {
			if i > 0 {
				form.Observer += "、"
			}
			form.Observer += v.User.Name
		}
		break
	case nil:
		break
	default:
		return dbReadError(c, err, template, &bind)
	}

	if c.Method() == "GET" {
		return c.Render(template, bind)
	}

	c.BodyParser(&form)
	// 處理時間
	if form.Date != "" && form.Time != "" {
		if _, err = time.Parse(dateFormat, form.Date); err != nil {
			return badRequest(c, "日期"+formatErrorMsg, template, &bind)
		}
		if form.Time != "" {
			form.Start, err = time.Parse(timeFormat, form.Date+form.Time)
			if err != nil {
				return badRequest(c, "開始時間"+formatErrorMsg, template, &bind)
			}
			form.Start = form.Start.Add(timeOffset).Local()
		}
	}

	err = form.Save()
	if err != nil {
		return dbWriteError(c, err, template, &bind)
	}

	bind["messages"] = []msgStruct{
		createMsg(infoMsgLevel, "儲存成功"),
	}

	return c.Render(template, bind)
}

// CourseBriefing 共同備課記錄表
func CourseBriefing(c *fiber.Ctx) error {
	return createCourseBriefingSH(c)
}
