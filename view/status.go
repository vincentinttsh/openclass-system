package view

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

func initStatusBind() {
	statusBind = map[string]string{
		"not_login":          "請先登入",
		"no_permission":      "沒有權限",
		"user_not_found":     "找不到使用者",
		"not_found":          "找不到資料",
		"is_login":           "您已經登入",
		"login":              "登入成功",
		"logout":             "登出成功",
		"already_registered": "您已經註冊過了",
		"create_success":     "新增成功",
		"create_fail":        "新增失敗",
		"update_success":     "更新成功",
		"update_fail":        "更新失敗",
		"delete_fail":        "刪除失敗",
		"delete_success":     "刪除成功",
		"attend_success":     "報到成功",
		"attend_fail":        "報到失敗",
		"reserve_success":    "預約成功",
		"reserve_fail":       "預約失敗",
	}
}

func statusBinding(c *fiber.Ctx, bind *fiber.Map) {
	status := c.Query("status", "")
	level := infoMsgLevel

	if status == "" {
		return
	}

	if strings.Contains(status, "fail") {
		level = errMsgLevel
	} else if status == "user_not_found" {
		level = errMsgLevel
	}

	(*bind)["messages"] = []msgStruct{
		createMsg(level, statusBind[status]),
	}
}
