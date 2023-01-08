package view

import "github.com/gofiber/fiber/v2"

// HomePage is a function that render home page
func HomePage(c *fiber.Ctx) error {
	return c.Render("home", fiber.Map{})
}
