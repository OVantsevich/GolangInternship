package handler

import (
	. "github.com/OVantsevich/GolangInternship/FMicroservice/internal/domain"
	. "github.com/OVantsevich/GolangInternship/FMicroservice/internal/service"
	"github.com/labstack/echo/v4"
	"net/http"
)

type UserHandler struct {
	es UserService
}

func NewUserHandler(es *UserService) *UserHandler {
	return &UserHandler{es: *es}
}

func (eh *UserHandler) Signup(c echo.Context) (err error) {
	User := &User{}
	if err = c.Bind(User); err != nil {
		return
	}

	if err = c.Validate(User); err != nil {
		return
	}

	if err = eh.es.CreateUser(c.Request().Context(), User); err != nil {
		return
	}

	return c.JSON(http.StatusCreated, User)
}

func (eh *UserHandler) Login(c echo.Context) (err error) {
	User := &User{}
	if err = c.Bind(User); err != nil {
		return
	}

	if User, err = eh.es.FindUser(c.Request().Context(), User.Name); err != nil {
		return
	}

	return c.JSON(http.StatusOK, User)
}

func (eh *UserHandler) Update(c echo.Context) (err error) {
	User := &User{}
	if err = c.Bind(User); err != nil {
		return
	}

	if err = c.Validate(User); err != nil {
		return
	}

	if err = eh.es.UpdateUser(c.Request().Context(), User); err != nil {
		return
	}

	return c.JSON(http.StatusOK, User)
}

func (eh *UserHandler) Delete(c echo.Context) (err error) {
	User := &User{}
	if err = c.Bind(User); err != nil {
		return
	}

	if err = eh.es.DeleteUser(c.Request().Context(), User.Name); err != nil {
		return
	}

	return c.JSON(http.StatusOK, User)
}
