package view

import (
	"vincentinttsh/openclass-system/model"

	"github.com/gofiber/fiber/v2"
)

// ListUserObservation 列出使用者的觀課紀錄
func ListUserObservation(c *fiber.Ctx) error {
	var userID = model.SQLBasePK(c.Locals("id").(float64))
	var template string = "observation/list"
	var bind fiber.Map = c.Locals("bind").(fiber.Map)
	var data []model.CourseReservation
	var reservations []model.CourseReservation
	var observations []model.CourseReservation
	var err error

	err = model.GetUserObserveCourses(&userID, &data)
	if err != nil {
		return dbReadError(c, err, template, &bind)
	}
	bind["reservations"] = &reservations
	bind["observations"] = &observations

	for _, reservation := range data {
		if reservation.Attended {
			reservations = append(reservations, reservation)
		} else {
			observations = append(observations, reservation)
		}
	}

	statusBinding(c, &bind)

	return c.Status(fiber.StatusOK).Render(template, bind)
}
