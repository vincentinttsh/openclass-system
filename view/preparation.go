package view

import (
	"strconv"
	"vincentinttsh/openclass-system/model"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func createCoursePreparationSH(c *fiber.Ctx) error {
	var classID model.SQLBasePK = 0
	var id uint64
	var courseData model.Course
	var form model.SHPreparation
	var template string = "preparation/formSH"
	var bind fiber.Map = c.Locals("bind").(fiber.Map)
	var ok bool
	var err error

	bind["method"] = "儲存"
	bind["permissions"] = "edit"
	bind["csrf_token"] = c.Locals("csrf_token")
	bind["form"] = &form

	if id, err = strconv.ParseUint(c.Params("id"), 10, 64); err != nil {
		return notFound(c)
	}
	classID = model.SQLBasePK(id)

	err = model.GetCourse(&classID, &courseData, true)
	switch err {
	case gorm.ErrRecordNotFound:
		return notFound(c)
	case nil:
		break
	default:
		return dbReadError(c, err, template, &bind)
	}

	err = model.GetSHPreparationByCourseID(&classID, &form)
	switch err {
	case gorm.ErrRecordNotFound:
		form.Course = courseData
		form.Subject = subjectChoice[courseData.User.Subject]
		form.Name = courseData.Name
		form.Date = courseData.Start.Format(dateFormat)
		form.StartTime = courseData.Start.Format(pureTimeFormat)
		form.EndTime = courseData.End.Format(pureTimeFormat)
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
	ok, err = dynamicTimeCheck(c, &form, template, &bind)
	if !ok {
		return err
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

// CoursePreparation 共同備課記錄表
func CoursePreparation(c *fiber.Ctx) error {
	return createCoursePreparationSH(c)
}
