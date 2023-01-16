package view

import (
	"vincentinttsh/openclass-system/model"

	"github.com/gofiber/fiber/v2"
)

type registerUser struct {
	Name       string `form:"username"`
	Email      string `form:"email"`
	Department string `form:"department" validate:"required"`
	Subject    string `form:"subject" validate:"required"`
}

var departmentChoice = []map[string]string{
	{"value": "sh", "text": "高中部"},
	{"value": "jh", "text": "國中部"},
}
var subjectChoice = []map[string]string{
	{"value": "chinese", "text": "國文"},
	{"value": "english", "text": "英文"},
	{"value": "math", "text": "數學"},
	{"value": "science", "text": "自然"},
	{"value": "social", "text": "社會"},
	{"value": "other", "text": "其他"},
}

// Register is a function that handles user registration
func Register(c *fiber.Ctx) error {
	var form registerUser
	var user model.User = c.Locals("user").(model.User)
	var updateFields model.User
	var err error

	if user.Subject != "" && user.Department != "" {
		return c.Redirect("/?status=already_registered")
	}

	if err = c.BodyParser(&form); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	errors := validateStruct(&form)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).Render("auth/register", fiber.Map{
			"messages": []msgStruct{
				createMsg(infoMsgLevel, "欄位不得為空。"),
			},
			"username":          form.Name,
			"email":             form.Email,
			"department":        form.Department,
			"subject":           form.Subject,
			"department_choice": departmentChoice,
			"subject_choice":    subjectChoice,
			"csrf_token":        c.Locals("csrf_token"),
		})
	}

	updateFields.Department = form.Department
	updateFields.Subject = form.Subject
	err = model.UpdateUser(&user, &updateFields)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).Render("auth/register", fiber.Map{
			"messages": []msgStruct{
				createMsg(errMsgLevel, err.Error()),
			},
			"username":          form.Name,
			"email":             form.Email,
			"department":        form.Department,
			"subject":           form.Subject,
			"department_choice": departmentChoice,
			"subject_choice":    subjectChoice,
			"csrf_token":        c.Locals("csrf_token"),
		})
	}

	return c.Redirect("/?status=login")
}
