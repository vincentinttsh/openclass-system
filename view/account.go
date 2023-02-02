package view

import (
	"time"
	"vincentinttsh/openclass-system/model"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Register is a function that handles user registration
func Register(c *fiber.Ctx) error {
	var template string = "auth/register"
	var form registerUser
	var user model.User
	var updateFields model.User
	var valid bool
	var err error
	var bind fiber.Map = c.Locals("bind").(fiber.Map)
	var userID = model.SQLBasePK(c.Locals("id").(float64))
	bind["department_choice"] = departmentChoice
	bind["subject_choice"] = subjectChoice
	bind["csrf_token"] = c.Locals("csrf_token")

	err = model.GetUserByID(&userID, &user)
	if err == gorm.ErrRecordNotFound {
		c.Cookie(setCookie("token", "", false, time.Now()))
		return c.RedirectToRoute("home", fiber.Map{
			"queries": map[string]string{
				"status": "user_not_found",
			},
		})
	}
	if err != nil {
		c.Cookie(setCookie("token", "", false, time.Now()))
		return serverError(c, err, template, &bind)
	}
	if user.Subject != "" && user.Department != "" {
		return c.RedirectToRoute("home", fiber.Map{
			"queries": map[string]string{
				"status": "already_registered",
			},
		})
	}
	if c.Method() == "GET" {
		return c.RedirectToRoute("authComplete", nil)
	}

	// readonly fields
	form.Username = user.Name
	form.Email = user.Email

	valid, err = formValidate(c, &bind, &form, &template)
	if !valid {
		return err
	}

	updateFields.Department = form.Department
	updateFields.Subject = form.Subject
	err = model.UpdateUser(&user, &updateFields)
	if err != nil {
		return serverError(c, err, template, &bind)
	}
	// resign jwt
	token, expire, err := signJWT(user)
	if err != nil {
		c.Cookie(setCookie("token", "", false, time.Now()))
		return serverError(c, err, template, &bind)
	}

	c.Cookie(setCookie("token", token, false, expire))

	return c.RedirectToRoute("home", fiber.Map{
		"queries": map[string]string{
			"status": "login",
		},
	})
}
