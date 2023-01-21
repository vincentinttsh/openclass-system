package view

import (
	"vincentinttsh/openclass-system/model"

	"github.com/gofiber/fiber/v2"
)

type registerUser struct {
	Username   string `form:"username" i18n:"姓名"`
	Email      string `form:"email" i18n:"Email address"`
	Department string `form:"department" i18n:"部門" validate:"required,oneof=sh jh"`
	Subject    string `form:"subject" i18n:"授課科目" validate:"required,oneof=chinese english math science social other"`
}

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

// Register is a function that handles user registration
func Register(c *fiber.Ctx) error {
	var username interface{} = c.Locals("username")
	var template string = "auth/register"
	var form registerUser
	var user model.User
	var updateFields model.User
	var valid bool
	var err error
	var bind fiber.Map = fiber.Map{
		"username":          username,
		"department_choice": departmentChoice,
		"subject_choice":    subjectChoice,
		"csrf_token":        c.Locals("csrf_token"),
	}

	user, err = model.GetUserByID(uint(c.Locals("id").(float64)))
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
