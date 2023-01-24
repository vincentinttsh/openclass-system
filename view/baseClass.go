package view

import (
	"time"
	"vincentinttsh/openclass-system/model"

	"github.com/gofiber/fiber/v2"
)

type baseClassStruct struct {
	Name      string `form:"className" validate:"required" i18n:"課程名稱"`
	Classroom string `form:"classroom" validate:"required" i18n:"上課教室"`
	Date      string `form:"date" validate:"required" i18n:"上課日期"`
	StartTime string `form:"startTime" validate:"required" i18n:"開始時間"`
	EndTime   string `form:"endTime" validate:"required" i18n:"結束時間"`
}

// CreateOpenClass create openclass class(POST, GET)
func CreateOpenClass(c *fiber.Ctx) error {
	var username interface{} = c.Locals("username")
	var userID = uint(c.Locals("id").(float64))
	var template string = "class/form"
	var form baseClassStruct
	var baseClass model.BaseClass
	var valid bool
	var err error
	var startTime time.Time
	var endTime time.Time
	var bind fiber.Map

	form = baseClassStruct{
		Date:      time.Now().Format("02/01/2006"),
		StartTime: time.Now().Format("15:04 PM"),
		EndTime:   time.Now().Add(time.Hour).Format("15:04 PM"),
	}
	bind = fiber.Map{
		"username":   username,
		"csrf_token": c.Locals("csrf_token"),
		"form":       form,
	}

	if c.Method() == fiber.MethodGet {
		return c.Status(fiber.StatusOK).Render(template, bind)
	}

	valid, err = formValidate(c, &bind, &form, &template)
	if !valid {
		return err
	}

	if _, err = time.Parse("02/01/2006", form.Date); err != nil {
		bind["messages"] = []msgStruct{
			createMsg(warnMsgLevel, "日期格式錯誤"),
		}
		return c.Status(fiber.StatusBadRequest).Render(template, bind)
	}
	if startTime, err = time.Parse("02/01/200603:04 PM", form.Date+form.StartTime); err != nil {
		bind["messages"] = []msgStruct{
			createMsg(warnMsgLevel, "開始時間格式錯誤"),
		}
		return c.Status(fiber.StatusBadRequest).Render(template, bind)
	}
	if endTime, err = time.Parse("02/01/200603:04 PM", form.Date+form.EndTime); err != nil {
		bind["messages"] = []msgStruct{
			createMsg(warnMsgLevel, "結束時間格式錯誤"),
		}
		return c.Status(fiber.StatusBadRequest).Render(template, bind)
	}

	// 轉換時區
	startTime = startTime.Add(time.Hour * -8)
	endTime = endTime.Add(time.Hour * -8)

	if endTime.Before(startTime) {
		bind["messages"] = []msgStruct{
			createMsg(warnMsgLevel, "結束時間不得早於開始時間"),
		}
		return c.Status(fiber.StatusBadRequest).Render(template, bind)
	}

	baseClass = model.BaseClass{
		Name:      form.Name,
		Classroom: form.Classroom,
		Start:     startTime.Local(),
		End:       endTime.Local(),
		TeacherID: userID,
	}

	if err = baseClass.Create(); err != nil {
		bind["messages"] = []msgStruct{
			createMsg(errMsgLevel, serverErrorMsg),
		}
		sugar.Errorln(err)
		return c.Status(fiber.StatusInternalServerError).Render(template, bind)
	}

	return c.Redirect("/?status=create_success")
}
