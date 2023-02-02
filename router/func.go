package router

import (
	"vincentinttsh/openclass-system/internal/jwt"

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
			return c.RedirectToRoute("login", fiber.Map{
				"queries": map[string]string{
					"status": "not_login",
					"next":   next,
				},
			})
		}

		return c.Next()
	}
}

func needRegistered() fiber.Handler {
	return func(c *fiber.Ctx) error {
		next := baseURL + c.OriginalURL()
		if c.Locals("id") == nil {
			return c.RedirectToRoute("login", fiber.Map{
				"queries": map[string]string{
					"status": "not_login",
					"next":   next,
				},
			})
		}
		if c.Locals("subject") == "" || c.Locals("department") == "" {
			return c.RedirectToRoute("authComplete", fiber.Map{})
		}

		return c.Next()
	}
}

func needSuperAdmin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		next := baseURL + c.OriginalURL()
		if c.Locals("id") == nil {
			return c.RedirectToRoute("login", fiber.Map{
				"queries": map[string]string{
					"status": "not_login",
					"next":   next,
				},
			})
		}
		if c.Locals("subject") == "" || c.Locals("department") == "" {
			return c.RedirectToRoute("authComplete", fiber.Map{})
		}
		if c.Locals("bind").(fiber.Map)["super_admin"] != true {
			return c.RedirectToRoute("home", fiber.Map{
				"queries": map[string]string{
					"status": "no_permission",
				},
			})
		}

		return c.Next()
	}
}
