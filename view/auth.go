package view

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/markbates/goth"
	"github.com/shareed2k/goth_fiber"
)

// Login is a function that authenticate user
func Login(c *fiber.Ctx) error {
	var err error
	var gothUser goth.User

	if gothUser, err = goth_fiber.CompleteUserAuth(c); err != nil {
		return c.Status(fiber.StatusConflict).SendString(fmt.Sprintf("%s", err))
	}

	return c.SendString(fmt.Sprintf("user email: %s\nuser picture: %s", gothUser.Email, gothUser.AvatarURL))
}

// Logout is a function that logout user
func Logout(c *fiber.Ctx) error {
	if err := goth_fiber.Logout(c); err != nil {
		return c.Status(fiber.StatusConflict).SendString(fmt.Sprintf("%s", err))
	}

	return c.SendString("logout")
}
