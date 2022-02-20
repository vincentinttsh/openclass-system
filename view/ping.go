package view

import (
	"github.com/gofiber/fiber/v2"
)

// Ping is a function that return the current application status
func Ping(c *fiber.Ctx) error {
	return c.Status(200).SendString("pong")
}
