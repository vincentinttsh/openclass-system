package view

import (
	"os"
	"time"
	"vincentinttsh/openclass-system/pkg/mode"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var googleOAuthConfig *oauth2.Config
var baseURL string
var signingKey []byte
var domain string

type profileStruct struct {
	ID     string `json:"id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	Locale string `json:"locale"`
	Domain string `json:"hd"`
}

func init() {
	googleOAuthConfig = &oauth2.Config{
		ClientID:     os.Getenv("OAUTH_KEY"),
		ClientSecret: os.Getenv("OAUTH_SECRET"),
		RedirectURL:  os.Getenv("OAUTH_CALLBACK_URL"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: google.Endpoint,
	}

	baseURL = os.Getenv("BASE_URL")
	signingKey = []byte(os.Getenv("JWT_SECRET"))
	domain = os.Getenv("DOMAIN")
}

func setCookie(name string, value string, session bool, expires time.Time) *fiber.Cookie {
	cookie := &fiber.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		Domain:   domain,
		Secure:   mode.Mode() == mode.ReleaseMode,
		HTTPOnly: true,
		SameSite: "Lax",
	}

	if session {
		cookie.SessionOnly = true
	}
	cookie.Expires = expires

	return cookie
}
