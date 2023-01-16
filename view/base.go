package view

import (
	"context"
	"os"
	"time"
	"vincentinttsh/openclass-system/pkg/mode"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/api/idtoken"
)

var baseURL string
var signingKey []byte
var domain string
var googleClientID string

var tokenValidator *idtoken.Validator
var validate *validator.Validate

const (
	infoMsgLevel = "info"
	errMsgLevel  = "error"
)

type errorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

type msgStruct map[string]string

func init() {
	var err error

	tokenValidator, err = idtoken.NewValidator(context.Background())
	if err != nil {
		panic(err)
	}

	baseURL = os.Getenv("BASE_URL")
	signingKey = []byte(os.Getenv("JWT_SECRET"))
	domain = os.Getenv("DOMAIN")
	googleClientID = os.Getenv("OAUTH_KEY")
	validate = validator.New()
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

func createMsg(level string, msg string) msgStruct {
	return msgStruct{
		"level": level,
		"msg":   msg,
	}
}

func validateStruct(data interface{}) []*errorResponse {
	var errors []*errorResponse
	err := validate.Struct(data)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element errorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}
