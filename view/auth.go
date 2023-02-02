package view

import (
	"context"
	"net/url"
	"strings"
	"time"

	"vincentinttsh/openclass-system/internal/jwt"
	"vincentinttsh/openclass-system/model"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/api/idtoken"
	"gorm.io/gorm"
)

func signJWT(user model.User) (t string, expire time.Time, err error) {
	expire = time.Now().Add(time.Hour * 72)
	claims := jwt.MapClaims{
		"id":                 user.ID,
		"name":               user.Name,
		"subject":            user.Subject,
		"department":         user.Department,
		"super_admin":        user.SuperAdmin,
		"admin":              user.Admin,
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
	var token string
	var payload *idtoken.Payload
	var profile map[string]interface{}
	var domain []string
	var user model.User
	var googleUser model.GoogleOauth
	var googleUID string

	csrfCookie := c.Cookies("g_csrf_token")
	csrfBody := c.FormValue("g_csrf_token")
	if csrfCookie == "" || csrfCookie != csrfBody {
		return c.Status(fiber.StatusUnauthorized).SendString("CSRF token mismatch")
	}

	token = c.FormValue("credential")
	payload, err = tokenValidator.Validate(context.Background(), token, googleClientID)
	if err != nil {
		sugar.Errorln(c)
		return c.Status(fiber.StatusUnauthorized).Render("base", fiber.Map{
			"messages": []msgStruct{
				createMsg(errMsgLevel, serverErrorMsg),
			},
		})
	}

	profile = payload.Claims
	domain = strings.Split(profile["hd"].(string), ".")
	googleUID = profile["sub"].(string)
	user, err = model.GetUserByGoogleID(&googleUID)
	if err == gorm.ErrRecordNotFound {
		googleUser = model.GoogleOauth{
			ID: profile["sub"].(string),
			User: model.User{
				Account: profile["email"].(string),
				Email:   profile["email"].(string),
				Name:    profile["name"].(string),
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
	return c.RedirectToRoute("authComplete", fiber.Map{})
}

// Complete is a function that redirect user to next page
func Complete(c *fiber.Ctx) error {
	var next *url.URL
	var template string = "auth/register"
	var user model.User
	var form registerUser
	var err error
	var bind = c.Locals("bind").(fiber.Map)
	var userID = model.SQLBasePK(c.Locals("id").(float64))
	bind["csrf_token"] = c.Locals("csrf_token")
	bind["department_choice"] = departmentChoice
	bind["subject_choice"] = subjectChoice

	err = model.GetUserByID(&userID, &user)
	if err == gorm.ErrRecordNotFound {
		c.Cookie(setCookie("token", "", false, time.Now()))
		return c.RedirectToRoute("home", fiber.Map{
			"queries": map[string]string{
				"status": "user_not_found",
			},
		})
	}
	if err != nil {
		return serverError(c, err, template, &bind)
	}
	if user.Subject == "" || user.Department == "" {
		form.Username = user.Name
		form.Email = user.Email
		bind["form"] = form
		bind["messages"] = []msgStruct{
			createMsg(infoMsgLevel, "請填寫以下資料，以完成註冊"),
		}
		return c.Status(fiber.StatusOK).Render(template, bind)
	}
	// Prevent open redirect vulnerability
	// Next URL must be same host
	next, err = url.ParseRequestURI(c.Cookies("redirect", ""))
	if err != nil || next.Scheme+"://"+next.Host != baseURL {
		errMsg := "no error"
		if err != nil {
			errMsg = err.Error()
		}
		sugar.Infow(
			"Redirect vulnerability attack",
			"IP", c.IP(),
			"Error", errMsg,
			"Next", next.String(),
			"username", user.Name,
			"email", user.Email,
		)
		return c.RedirectToRoute("home", fiber.Map{
			"queries": map[string]string{
				"status": "login",
			},
		})
	}
	return c.Redirect(next.String())
}

// Logout is a function that logout user
func Logout(c *fiber.Ctx) error {
	c.Cookie(setCookie("token", "", false, time.Now()))
	return c.RedirectToRoute("home", fiber.Map{
		"queries": map[string]string{
			"status": "logout",
		},
	})
}

// LoginPage is a function that render login page
func LoginPage(c *fiber.Ctx) error {
	var bind fiber.Map = c.Locals("bind").(fiber.Map)
	var next = c.Query("next", baseURL+"?status=login")

	if c.Locals("id") != nil {
		return c.RedirectToRoute("home", fiber.Map{
			"queries": map[string]string{
				"status": "is_login",
			},
		})
	}
	c.Cookie(setCookie("redirect", next, true, time.Now()))
	statusBinding(c, &bind)

	return c.Status(fiber.StatusOK).Render("auth/login", bind)
}
