package view

import (
	"github.com/gofiber/fiber/v2"
)

// HomePage is a function that render home page
func HomePage(c *fiber.Ctx) error {
	username := c.Locals("username")
	return c.Render("home", fiber.Map{
		"username": username,
	})
}
