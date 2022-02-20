package view

import (
	"fmt"
	"vincentinttsh/openclass-system/config"
	"vincentinttsh/openclass-system/pkg/mode"

	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/fiber/v2/utils"
)

var store *session.Store

// InitFunc is a function that initializes the view
func InitFunc(config *config.Config) {
	store = session.New(session.Config{
		Storage:        config.SessionStorage,
		CookieDomain:   config.Domain,
		CookieSecure:   mode.Mode() == mode.ReleaseMode,
		CookieHTTPOnly: true,
		CookieSameSite: "Lax",
		KeyGenerator:   utils.UUIDv4,
	})
}

func print(value interface{}) {
	fmt.Printf("%+v\n", value)
}
