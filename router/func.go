package router

import (
	"vincentinttsh/openclass-system/pkg/jwt"

	"github.com/gofiber/fiber/v2"
)

func jwtVerify() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var tokenString string
		tokenString = c.Cookies("token", "")
		if tokenString == "" {
			return c.Next()
		}

		claims, err := jwt.Verify(tokenString)
		if err != nil {
			return c.Next()
		}

		c.Locals("username", claims["name"])
		c.Locals("id", claims["id"])

		return c.Next()
	}
}
