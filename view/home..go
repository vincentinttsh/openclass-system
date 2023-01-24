package view

import (
	"fmt"
	"vincentinttsh/openclass-system/model"

	"github.com/gofiber/fiber/v2"
)

type classBind struct {
	ClassID   uint   `json:"class_id"`
	ClassName string `json:"class_name"`
	Classroom string `json:"classroom"`
	Date      string `json:"date"`
	Duration  string `json:"duration"`
	Teacher   string `json:"teacher"`
	Calendar  string `json:"calendar"`
}

// HomePage is a function that render home page
func HomePage(c *fiber.Ctx) error {
	var username interface{} = c.Locals("username")
	var data []model.BaseClass
	var classes []classBind
	var err error
	var bind fiber.Map = fiber.Map{
		"username": username,
	}

	err = model.GetAllClass(&data)
	if err != nil {
		sugar.Errorw("Get all class error", "error", err)
		bind["messages"] = []msgStruct{
			createMsg(errMsgLevel, "取得課程資料時發生錯誤"),
		}
	}

	classes = make([]classBind, len(data))
	var calendar string = "http://www.google.com/calendar/event?action=TEMPLATE&text=%s公開授課（"
	calendar += "%s)&dates=%s/%s&details=課程名稱：%s"
	for i, v := range data {
		classes[i] = classBind{
			ClassID:   v.ID,
			ClassName: v.Name,
			Classroom: v.Classroom,
			Date:      v.Start.Format("2006年01月02日"),
			Duration:  v.Start.Format("15:04") + " ~ " + v.End.Format("15:04"),
			Teacher:   v.Teacher.Name,
			Calendar: fmt.Sprintf(calendar,
				departmentChoice[v.Teacher.Department],
				v.Teacher.Name,
				v.Start.Format("20060102T150405"),
				v.End.Format("20060102T150405"),
				v.Name,
			) + "%0A" + fmt.Sprintf("授課老師：%s&location=%s&trp=false", v.Teacher.Name, v.Classroom),
		}
	}

	bind["classes"] = &classes

	switch c.Query("status", "") {
	case "login":
		bind["messages"] = []msgStruct{
			createMsg(infoMsgLevel, "登入成功"),
		}
	case "logout":
		bind["messages"] = []msgStruct{
			createMsg(infoMsgLevel, "登出成功"),
		}
	case "already_registered":
		bind["messages"] = []msgStruct{
			createMsg(infoMsgLevel, "您已經註冊過了"),
		}
	case "create_success":
		bind["messages"] = []msgStruct{
			createMsg(infoMsgLevel, "新增成功"),
		}
	}

	return c.Status(fiber.StatusOK).Render("home", bind)
}
