package view

import (
	"context"
	"os"
	"reflect"
	"time"
	"vincentinttsh/openclass-system/model"
	"vincentinttsh/openclass-system/pkg/mode"
	"vincentinttsh/openclass-system/pkg/tool"

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

var timeOffset time.Duration

// ErrorResponse is a struct for error response
type ErrorResponse struct {
	Field string
	Tag   string
}

type msgStruct struct {
	Level string
	Msg   string
}

func init() {
	var err error

	tokenValidator, err = idtoken.NewValidator(context.Background())
	if err != nil {
		panic(err)
	}

	time.LoadLocation(os.Getenv("TIMEZONE"))
	baseURL = os.Getenv("BASE_URL")
	signingKey = []byte(os.Getenv("JWT_SECRET"))
	domain = os.Getenv("DOMAIN")
	googleClientID = os.Getenv("OAUTH_KEY")
	validate = validator.New()
	_, tmp := time.Now().Zone()
	timeOffset = -1 * time.Duration(tmp) * time.Second
	initLogger()
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
		Level: level,
		Msg:   msg,
	}
}

func validateStruct(data interface{}) []*ErrorResponse {
	var errors []*ErrorResponse
	var err error

	err = validate.Struct(data)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.Field = err.Field()
			element.Tag = err.Tag()
			errors = append(errors, &element)
		}
	}
	return errors
}

func formValidate(c *fiber.Ctx, bind *fiber.Map, form interface{}, template *string) (bool, error) {

	// form is reference, so can bind here
	(*bind)["form"] = form

	if err := c.BodyParser(form); err != nil {
		(*bind)["messages"] = []msgStruct{
			createMsg(errMsgLevel, serverErrorMsg),
		}
		sugar.Errorln(c.Body(), err)
		return false, c.Status(fiber.StatusInternalServerError).Render(*template, *bind)
	}

	errors := validateStruct(form)
	if errors != nil {
		var msg []msgStruct
		for _, err := range errors {
			tool.Print(err)
			msg = append(msg, createMsg(warnMsgLevel, formErrorHandler(err, form)))
		}
		(*bind)["messages"] = msg
		return false, c.Status(fiber.StatusBadRequest).Render(*template, *bind)
	}

	return true, nil
}

func formErrorHandler(er *ErrorResponse, form interface{}) string {
	var field, match = reflect.TypeOf(form).Elem().FieldByName(er.Field)
	if !match {
		sugar.Errorln("field not found: " + er.Field)
		return "欄位錯誤"
	}
	var msg string = field.Tag.Get("i18n")
	switch er.Tag {
	case "required":
		msg += ": 為必填欄位"
	case "oneof":
		msg += ": 輸入值錯誤"
	}
	tool.Print(msg)
	return msg
}

type baseTimeStruct interface {
	model.SHPreparation | model.SHDesign
}

func dynamicTimeCheck(
	c *fiber.Ctx, form model.DurationBaseInterface, template string, bind *fiber.Map,
) (bool, error) {
	date, startTime, endTime := form.GetTimeString()
	var start time.Time
	var end time.Time

	var err error
	if date != "" {
		if _, err = time.Parse(dateFormat, date); err != nil {
			return false, badRequest(c, "日期"+formatErrorMsg, template, bind)
		}
		if startTime != "" {
			start, err = time.Parse(timeFormat, date+startTime)
			if err != nil {
				return false, badRequest(c, "開始時間"+formatErrorMsg, template, bind)
			}
			start = start.Add(timeOffset).Local()
			form.SetStartTime(start)
		}
		if endTime != "" {
			end, err = time.Parse(timeFormat, date+endTime)
			if err != nil {
				return false, badRequest(c, "結束時間"+formatErrorMsg, template, bind)
			}
			end = end.Add(timeOffset).Local()
			form.SetEndTime(end)
		}
		if startTime != "" && endTime != "" && end.Before(start) {
			return false, badRequest(c, "結束時間不得早於開始時間", template, bind)
		}
	}
	return true, nil
}
