package router

import (
	"os"
	"vincentinttsh/openclass-system/internal/mode"
	"vincentinttsh/openclass-system/view"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/storage/redis"
)

var signingKey []byte
var baseURL string

// SetupRouter is a function that sets up the router
func SetupRouter(app *fiber.App) {
	baseURL = os.Getenv("BASE_URL")
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
	app.Use(csrf.New(csrfConfig))
	app.Get("/login", view.LoginPage).Name("login")
	app.Get("/logout", view.Logout)
	app.Get("/", view.HomePage).Name("home")

	// need login
	app.Get("/auth/complete", needLogin(), view.Complete).Name("authComplete")
	app.Get("/register", needLogin(), view.Register)
	app.Post("/register", needLogin(), view.Register)

	// need registered (with CSRF)
	app.Get("/class/create", needRegistered(), view.CreateOpenClass)
	app.Post("/class/create", needRegistered(), view.CreateOpenClass)
	app.Get("/class/:id", needRegistered(), view.GetOfModifyOpenClass)
	app.Post("/class/:id", needRegistered(), view.GetOfModifyOpenClass)
	app.Get("/class/:id/design", needRegistered(), view.CourseDesign)
	app.Post("/class/:id/design", needRegistered(), view.CourseDesign)
	app.Get("/class/:id/preparation", needRegistered(), view.CoursePreparation)
	app.Post("/class/:id/preparation", needRegistered(), view.CoursePreparation)
	app.Get("/class/:id/briefing", needRegistered(), view.CourseBriefing)
	app.Post("/class/:id/briefing", needRegistered(), view.CourseBriefing)

	// need registered  (without CSRF)
	app.Get("/class/:id/participation", needRegistered(), view.AttendClass)
	app.Get("/class/:id/reservation", needRegistered(), view.ReserveClass)
	app.Get("/my/class", needRegistered(), view.ListUserOpenClass).Name("my_class")
	app.Get("/my/observation", needRegistered(), view.ListUserObservation).Name("my_observation")
	app.Get("/metrics", needSuperAdmin(), monitor.New(monitor.Config{Title: "MyService Metrics Page"}))
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).Render("error/404", fiber.Map{
			"status": fiber.StatusNotFound,
			"title":  "找不到頁面",
		})
	})
}
