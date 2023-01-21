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
	// auth
	app.Post("/auth/callback/:provider", view.Login)

	app.Use(jwtVerify())
	app.Get("/metrics", monitor.New(monitor.Config{Title: "MyService Metrics Page"}))
	app.Get("/login", view.LoginPage)
	app.Get("/logout", view.Logout)
	// home
	app.Get("/", view.HomePage)

	// need login
	needLogin := app.Group("", needLogin())
	needCSRF := needLogin.Group("", csrf.New(csrfConfig))
	needCSRF.Get("/auth/complete", view.Complete)
	needCSRF.Get("/register", view.Register)
	needCSRF.Post("/register", view.Register)

	// need registered
	needRegistered := needCSRF.Group("", needRegistered())
	needRegistered.Get("/class/create", view.CreateOpenClass)
	needRegistered.Post("/class/create", view.CreateOpenClass)
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).SendString("Sorry can't find that!")
	})
}
