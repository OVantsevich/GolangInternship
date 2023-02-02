package handler

import (
	"github.com/labstack/echo/v4"
)

type User interface {
	Signup(c echo.Context) (err error)
	Login(c echo.Context) (err error)
	Refresh(c echo.Context) (err error)
	Update(c echo.Context) (err error)
	Delete(c echo.Context) (err error)

	UserByLogin(c echo.Context) (err error)
}
