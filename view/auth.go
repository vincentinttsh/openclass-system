package view

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"vincentinttsh/openclass-system/model"
	"vincentinttsh/openclass-system/pkg/jwt"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

// BeginAuthHandler is a function that redirect user to google oauth page
func BeginAuthHandler(c *fiber.Ctx) error {
	c.Cookie(setCookie("redirect", c.Query("redirect"), true, time.Now()))

	return c.Redirect(googleOAuthConfig.AuthCodeURL("state"))
}

// googleOauth is a function that authenticate user with google oauth
func googleOauth(c *fiber.Ctx) (profile profileStruct, err error) {
	var token *oauth2.Token
	var client *http.Client
	var res *http.Response

	if token, err = googleOAuthConfig.Exchange(oauth2.NoContext, c.Query("code")); err != nil {
		err = errors.New("token exchange: " + err.Error())
		return
	}
	client = googleOAuthConfig.Client(oauth2.NoContext, token)
	if res, err = client.Get("https://www.googleapis.com/oauth2/v2/userinfo"); err != nil {
		err = errors.New("get info" + err.Error())
		return
	}
	defer res.Body.Close()

	bytes, _ := io.ReadAll(res.Body)
	if err = json.Unmarshal(bytes, &profile); err != nil {
		err = errors.New("unmarshal: " + err.Error())
	}

	return
}

func signJWT(user model.User) (t string, expire time.Time, err error) {
	expire = time.Now().Add(time.Hour * 72)
	claims := jwt.MapClaims{
		"id":                 user.ID,
		"name":               user.Name,
		"super_admin":        user.SuperAdmin,
		"admin":              user.Admin,
		"organization_id":    user.OrganizationID,
		"organization_abbr":  user.Organization.Abbr,
		"organization_level": user.Organization.Level,
		"exp":                expire.Unix(),
	}

	t, err = jwt.Sign(claims)

	return
}

// Login is a function that authenticate user
func Login(c *fiber.Ctx) error {
	var err error
	var profile profileStruct

	if profile, err = googleOauth(c); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	var domain []string
	var user model.User
	var googleUser model.GoogleOauth
	domain = strings.Split(profile.Domain, ".")
	user, err = model.GetUserByGoogleID(profile.ID)
	if err == gorm.ErrRecordNotFound {
		googleUser = model.GoogleOauth{
			ID: profile.ID,
			User: model.User{
				Account: profile.Email,
				Email:   profile.Email,
				Name:    profile.Name,
				Locale:  profile.Locale,
				Organization: model.Organization{
					Level: domain[1],
					Abbr:  domain[0],
				},
			},
		}

		if err = model.CreateUserFromGoogle(&googleUser); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		user = googleUser.User
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	token, expire, err := signJWT(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	c.Cookie(setCookie("token", token, false, expire))

	return c.Redirect("/auth/complete")
}

// Complete is a function that redirect user to next page
func Complete(c *fiber.Ctx) error {
	// Prevent open redirect vulnerability
	// Next URL must be same host
	next, err := url.ParseRequestURI(c.Cookies("redirect", baseURL))
	if err != nil || next.Host != c.Hostname() {
		next = &url.URL{Path: baseURL}
	}

	return c.Redirect(next.String())
}

// Logout is a function that logout user
func Logout(c *fiber.Ctx) error {
	c.Cookie(setCookie("token", "", false, time.Now()))
	return c.Redirect("/")
}

// LoginPage is a function that render login page
func LoginPage(c *fiber.Ctx) error {
	return c.Render("auth/login", fiber.Map{})
}
