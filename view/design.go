package view

import (
	"strconv"
	"time"
	"vincentinttsh/openclass-system/model"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func createDesignSH(c *fiber.Ctx) error {
	var classID model.SQLBasePK = 0
	var id uint64
	var template string = "design/formSH"
	var courseData model.Course
	var form model.SHDesign
	var bind fiber.Map = c.Locals("bind").(fiber.Map)
	var oldDetails []model.SHDesignDetail
	var err error

	bind["method"] = "儲存"
	bind["permissions"] = "edit"
	bind["csrf_token"] = c.Locals("csrf_token")

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

	// 取得高中課程教學活動設計表, 若不存在則填入課程資料
	err = model.GetSHDesignByCourseID(&classID, &form)
	switch err {
	case gorm.ErrRecordNotFound:
		form.Course = courseData
		form.Name = courseData.Name
		form.Date = courseData.Start.Format(dateFormat)
		form.StartTime = courseData.Start.Format(pureTimeFormat)
		form.EndTime = courseData.End.Format(pureTimeFormat)
		form.Details = make([]model.SHDesignDetail, 1)
		break
	case nil:
		break
	default:
		return dbReadError(c, err, template, &bind)
	}

	bind["form"] = &form

	if c.Method() == "GET" {
		return c.Render(template, bind)
	}

	oldDetails = form.Details
	form.Details = make([]model.SHDesignDetail, 0)

	c.BodyParser(&form)
	// 處理時間
	if form.Date != "" {
		if _, err = time.Parse(dateFormat, form.Date); err != nil {
			return badRequest(c, "日期"+formatErrorMsg, template, &bind)
		}
		if form.StartTime != "" {
			form.Start, err = time.Parse(timeFormat, form.Date+form.StartTime)
			if err != nil {
				return badRequest(c, "開始時間"+formatErrorMsg, template, &bind)
			}
			form.Start = form.Start.Add(timeOffset).Local()
		}
		if form.EndTime != "" {
			form.End, err = time.Parse(timeFormat, form.Date+form.EndTime)
			if err != nil {
				return badRequest(c, "結束時間"+formatErrorMsg, template, &bind)
			}
			form.End = form.End.Add(timeOffset).Local()
		}
		if form.StartTime != "" && form.EndTime != "" && form.End.Before(form.Start) {
			return badRequest(c, "結束時間不得早於開始時間", template, &bind)
		}
	}

	err = form.Update()
	if err != nil {
		return dbWriteError(c, err, template, &bind)
	}

	for i := range form.Details {
		if i < len(oldDetails) {
			form.Details[i].ID = oldDetails[i].ID
		}
		form.Details[i].SHDesignID = form.ID
		err = form.Details[i].Update()
		if err != nil {
			return dbWriteError(c, err, template, &bind)
		}
	}
	for i := len(form.Details); i < len(oldDetails); i++ {
		err = oldDetails[i].Delete()
		if err != nil {
			return dbWriteError(c, err, template, &bind)
		}
	}

	bind["messages"] = []msgStruct{
		createMsg(infoMsgLevel, "儲存成功"),
	}

	return c.Render(template, bind)
}

// GetCourseDesign is the function to get or create the design of the class
func GetCourseDesign(c *fiber.Ctx) error {
	return createDesignSH(c)
}
