package router

import (
	"os"
	"vincentinttsh/openclass-system/config"
	"vincentinttsh/openclass-system/pkg/mode"
	"vincentinttsh/openclass-system/view"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"
	"github.com/shareed2k/goth_fiber"
)

// SetupRouter is a function that sets up the router
func SetupRouter(app *fiber.App, config *config.Config) {

	app.Use(favicon.New(favicon.Config{
		File: "./web/static/favicon.ico",
	}))

	app.Use(logger.New(logger.Config{
		Format:     "[${time}] ${status} - ${latency} ${method} ${path}\n",
		TimeFormat: "2006-Jan-02 15:04:05",
		TimeZone:   "Asia/Taipei",
	}))

	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))

	goth_fiber.SessionStore = session.New(session.Config{
		CookieDomain:   config.Domain,
		CookieSecure:   mode.Mode() == mode.ReleaseMode,
		CookieHTTPOnly: true,
		CookieSameSite: "Lax",
		KeyGenerator:   utils.UUIDv4,
		KeyLookup:      "cookie:auth_session",
	})

	goth.UseProviders(
		google.New(os.Getenv("OAUTH_KEY"), os.Getenv("OAUTH_SECRET"), os.Getenv("OAUTH_CALLBACK_URL")),
	)

	// auth
	app.Get("/auth/:provider", goth_fiber.BeginAuthHandler)
	app.Get("/auth/callback/:provider", view.Login)
	app.Get("/logout", view.Logout)

	// Ping
	app.Get("/ping", view.Ping)
}
