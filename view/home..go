package view

import (
	"vincentinttsh/openclass-system/model"

	"github.com/gofiber/fiber/v2"
)

// HomePage is a function that render home page
func HomePage(c *fiber.Ctx) error {
	var courses []model.Course
	var err error
	var bind fiber.Map = c.Locals("bind").(fiber.Map)

	err = model.GetAllCourses(&courses)
	if err != nil {
		sugar.Errorw("Get all class error", "error", err)
		bind["messages"] = []msgStruct{
			createMsg(errMsgLevel, "取得課程資料時發生錯誤"),
		}
	}

	bind["courses"] = &courses
	statusBinding(c, &bind)

	return c.Status(fiber.StatusOK).Render("home", bind)
}
