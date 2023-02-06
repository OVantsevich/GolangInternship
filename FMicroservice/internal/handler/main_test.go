package handler

import (
	"GolangInternship/FMicroservice/internal/model"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"net/http"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return nil
}

var e *echo.Echo

var testValidData = []model.User{
	{
		Name:     `NAME`,
		Age:      5,
		Login:    `CreateLOGIN1`,
		Email:    `LOGIN1@gmail.com`,
		Token:    `validToken`,
		Password: `strongPassword`,
	},
	{
		Name:     `NAME`,
		Age:      5,
		Login:    `CreateLOGIN2`,
		Email:    `LOGIN2@gmail.com`,
		Token:    `validToken2`,
		Password: `PASSWORD123456789`,
	},
}
var testNoValidData = []model.User{
	{
		Name:     `nameEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEE`,
		Age:      22222,
		Login:    `LOGIN2`,
		Email:    `LOGIN2@gmail.com`,
		Token:    `dafrawerfaegfaegae`,
		Password: `weak`,
	},
	{
		Name:     `NAME`,
		Age:      2,
		Login:    `LOGIN1AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA`,
		Email:    `LOGIN1@gmail.com`,
		Token:    `argawegfafawfew`,
		Password: `lalala`,
	},
}

func testInit() {
	if e == nil {
		e = echo.New()
		e.Validator = &CustomValidator{validator: validator.New()}
	}
}
