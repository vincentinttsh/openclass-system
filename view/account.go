package view

import (
	"vincentinttsh/openclass-system/model"

	"github.com/gofiber/fiber/v2"
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
	bind["department_choice"] = departmentChoice
	bind["subject_choice"] = subjectChoice
	bind["csrf_token"] = c.Locals("csrf_token")

	user, err = model.GetUserByID(model.SQLBasePK(c.Locals("id").(float64)))
	if err != nil {
		return c.Redirect("/login?status=notfound", fiber.StatusFound)
	}
	if user.Subject != "" && user.Department != "" {
		return c.Redirect("/?status=already_registered")
	}
	if c.Method() == "GET" {
		return c.Redirect("/auth/complete")
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
		sugar.Errorln(err)
		bind["messages"] = []msgStruct{
			createMsg(errMsgLevel, serverErrorMsg),
		}
		return c.Status(fiber.StatusInternalServerError).Render(template, bind)
	}
	// resign jwt
	token, expire, err := signJWT(user)
	if err != nil {
		sugar.Errorln(err)
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	c.Cookie(setCookie("token", token, false, expire))

	return c.Redirect("/?status=login")
}
