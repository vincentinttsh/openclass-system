package view

import (
	"github.com/gofiber/fiber/v2"
)

const (
	readDBdataErrorMsg  string = "取得資料時發生錯誤"
	writeDBdataErrorMsg string = "寫入資料時發生錯誤"
	serverErrorMsg      string = "伺服器錯誤，請連絡管理員"
	formatErrorMsg      string = "格式錯誤"
)

func notFound(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotFound).Render("error/404", fiber.Map{
		"status": fiber.StatusNotFound,
		"title":  "找不到頁面",
	})
}

func badRequest(c *fiber.Ctx, msg string, template string, bind *fiber.Map) error {
	(*bind)["messages"] = []msgStruct{
		createMsg(warnMsgLevel, msg),
	}
	return c.Status(fiber.StatusBadRequest).Render(template, *bind)
}

func serverError(c *fiber.Ctx, err error, template string, bind *fiber.Map) error {
	sugar.Errorln(err)
	(*bind)["messages"] = []msgStruct{
		createMsg(errMsgLevel, serverErrorMsg),
	}
	return c.Status(fiber.StatusInternalServerError).Render(template, *bind)
}

func dbReadError(c *fiber.Ctx, err error, template string, bind *fiber.Map) error {
	sugar.Errorln(err)
	(*bind)["messages"] = []msgStruct{
		createMsg(errMsgLevel, readDBdataErrorMsg),
	}
	return c.Status(fiber.StatusInternalServerError).Render(template, *bind)
}

func dbWriteError(c *fiber.Ctx, err error, template string, bind *fiber.Map) error {
	sugar.Errorln(err)
	(*bind)["messages"] = []msgStruct{
		createMsg(errMsgLevel, writeDBdataErrorMsg),
	}
	return c.Status(fiber.StatusInternalServerError).Render(template, *bind)
}
