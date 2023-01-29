package router

import (
	"vincentinttsh/openclass-system/pkg/jwt"

	"github.com/gofiber/fiber/v2"
)

func jwtVerify() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var tokenString string
		tokenString = c.Cookies("token", "")
		c.Locals("bind", fiber.Map{})
		if tokenString == "" {
			return c.Next()
		}

		claims, err := jwt.Verify(tokenString)
		if err != nil {
			return c.Next()
		}

		c.Locals("id", claims["id"])
		c.Locals("subject", claims["subject"])
		c.Locals("department", claims["department"])
		c.Locals("bind", fiber.Map{
			"super_admin": claims["super_admin"],
			"admin":       claims["admin"],
			"username":    claims["name"],
		})

		return c.Next()
	}
}

func needLogin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		next := baseURL + c.OriginalURL()
		if c.Locals("id") == nil {
			return c.Redirect("/login?status=notfound&next=" + next)
		}

		return c.Next()
	}
}

func needRegistered() fiber.Handler {
	return func(c *fiber.Ctx) error {
		next := baseURL + c.OriginalURL()
		if c.Locals("id") == nil {
			return c.Redirect("/login?status=notfound&next=" + next)
		}
		if c.Locals("subject") == "" || c.Locals("department") == "" {
			return c.Redirect("/auth/complete")
		}

		return c.Next()
	}
}

func needSuperAdmin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		next := baseURL + c.OriginalURL()
		if c.Locals("id") == nil {
			return c.Redirect("/login?status=notfound&next=" + next)
		}
		if c.Locals("subject") == "" || c.Locals("department") == "" {
			return c.Redirect("/auth/complete")
		}
		if c.Locals("bind").(fiber.Map)["super_admin"] != true {
			return c.Redirect("/?status=permission")
		}

		return c.Next()
	}
}
