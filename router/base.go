package router

import (
	"os"
	"vincentinttsh/openclass-system/pkg/mode"
	"vincentinttsh/openclass-system/view"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/storage/redis"
)

var signingKey []byte

// SetupRouter is a function that sets up the router
func SetupRouter(app *fiber.App) {
	signingKey = []byte(os.Getenv("JWT_SECRET"))
	csrfConfig := csrf.Config{
		KeyLookup: "form:csrf_token",
		Storage: redis.New(redis.Config{
			URL: os.Getenv("REDIS_URL"),
		}),
		CookieSessionOnly: true,
		CookieHTTPOnly:    true,
		CookieDomain:      os.Getenv("DOMAIN"),
		CookieSecure:      mode.Mode() == mode.ReleaseMode,
		CookieSameSite:    "Strict",
		ContextKey:        "csrf_token",
	}

	// Ping
	app.Get("/ping", view.Ping)

	app.Use(func(c *fiber.Ctx) error {
		c.Append(fiber.HeaderReferrerPolicy, "no-referrer-when-downgrade")
		return c.Next()
	})
	app.Post("/auth/callback/:provider", view.Login)

	app.Use(jwtVerify())
	app.Get("/metrics", monitor.New(monitor.Config{Title: "MyService Metrics Page"}))
	app.Get("/login", view.LoginPage)
	app.Get("/logout", view.Logout)
	app.Get("/", view.HomePage)

	// need login
	needLoginPath := app.Group("", needLogin())
	needCSRFPath := needLoginPath.Group("", csrf.New(csrfConfig))
	needCSRFPath.Get("/auth/complete", view.Complete)
	needCSRFPath.Get("/register", view.Register)
	needCSRFPath.Post("/register", view.Register)

	// need registered (with CSRF)
	needRegisteredPath := needCSRFPath.Group("", needRegistered())
	needRegisteredPath.Get("/class/create", view.CreateOpenClass)
	needRegisteredPath.Post("/class/create", view.CreateOpenClass)
	needRegisteredPath.Get("/class/:id", view.GetOfModifyOpenClass)
	needRegisteredPath.Post("/class/:id", view.GetOfModifyOpenClass)

	// need registered  (without CSRF)
	needRegisteredPath = needLoginPath.Group("", needRegistered())
	needRegisteredPath.Get("/my/class", view.ListUserOpenClass)
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).Render("error/404", fiber.Map{
			"status": fiber.StatusNotFound,
			"title":  "找不到頁面",
		})
	})
}
