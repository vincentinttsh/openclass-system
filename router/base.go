package router

import (
	"os"
	"vincentinttsh/openclass-system/view"

	"github.com/gofiber/fiber/v2"
)

var signingKey []byte

// SetupRouter is a function that sets up the router
func SetupRouter(app *fiber.App) {
	signingKey = []byte(os.Getenv("JWT_SECRET"))
	// Ping
	app.Get("/ping", view.Ping)

	// auth
	app.Get("/auth/provider/:provider", view.BeginAuthHandler)
	app.Get("/auth/callback/:provider", view.Login)
	app.Get("/auth/complete", view.Complete)

	app.Use(jwtVerify())
	app.Get("/login", view.LoginPage)
	app.Get("/logout", view.Logout)

	// home
	app.Get("/", view.HomePage)
}
