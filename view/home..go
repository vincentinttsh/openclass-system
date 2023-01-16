package view

import (
	"github.com/gofiber/fiber/v2"
)

// HomePage is a function that render home page
func HomePage(c *fiber.Ctx) error {
	var username interface{} = c.Locals("username")
	var status string = c.Query("status", "")

	var bind fiber.Map = fiber.Map{
		"username": username,
	}

	switch status {
	case "login":
		bind["messages"] = []msgStruct{
			createMsg(infoMsgLevel, "登入成功。"),
		}
	case "logout":
		bind["messages"] = []msgStruct{
			createMsg(infoMsgLevel, "登出成功。"),
		}
	case "already_registered":
		bind["messages"] = []msgStruct{
			createMsg(infoMsgLevel, "您已經註冊過了。"),
		}
	}

	return c.Render("home", bind)
}
