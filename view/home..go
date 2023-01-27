package view

import (
	"fmt"
	"vincentinttsh/openclass-system/model"

	"github.com/gofiber/fiber/v2"
)

type courseBind struct {
	ClassID   model.SQLBasePK `json:"class_id"`
	ClassName string          `json:"class_name"`
	Classroom string          `json:"classroom"`
	Date      string          `json:"date"`
	Duration  string          `json:"duration"`
	Teacher   string          `json:"teacher"`
	Calendar  string          `json:"calendar"`
	Passwd    string          `json:"passwd"`
}

// HomePage is a function that render home page
func HomePage(c *fiber.Ctx) error {
	var data []model.Course
	var courses []courseBind
	var err error
	var bind fiber.Map = c.Locals("bind").(fiber.Map)

	err = model.GetAllCourses(&data)
	if err != nil {
		sugar.Errorw("Get all class error", "error", err)
		bind["messages"] = []msgStruct{
			createMsg(errMsgLevel, "取得課程資料時發生錯誤"),
		}
	}

	courses = make([]courseBind, len(data))
	var calendar string = "http://www.google.com/calendar/event?action=TEMPLATE&text=%s公開授課（"
	calendar += "%s)&dates=%s/%s&details=課程名稱：%s"
	for i, v := range data {
		courses[i] = courseBind{
			ClassID:   v.ID,
			ClassName: v.Name,
			Classroom: v.Classroom,
			Date:      v.Start.Format("2006年01月02日"),
			Duration:  v.Start.Format("15:04") + " ~ " + v.End.Format("15:04"),
			Teacher:   v.User.Name,
			Calendar: fmt.Sprintf(calendar,
				departmentChoice[v.User.Department],
				v.User.Name,
				v.Start.Format("20060102T150405"),
				v.End.Format("20060102T150405"),
				v.Name,
			) + "%0A" + fmt.Sprintf("授課老師：%s&location=%s&trp=false", v.User.Name, v.Classroom),
		}
	}

	bind["courses"] = &courses

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
	case "update_success":
		bind["messages"] = []msgStruct{
			createMsg(infoMsgLevel, "更新成功"),
		}
	}

	return c.Status(fiber.StatusOK).Render("home", bind)
}
