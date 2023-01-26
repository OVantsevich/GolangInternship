package handler

import (
	. "github.com/OVantsevich/GolangInternship/FMicroservice/internal/domain"
	. "github.com/OVantsevich/GolangInternship/FMicroservice/internal/service"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"net/http"
)

type (
	UserHandler struct {
		es UserService
	}
)

type TokenResponse struct {
	AccessToken  string `json:"access"`
	RefreshToken string `json:"refresh"`
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

	var accessToken, refreshToken string
	if accessToken, refreshToken, err = eh.es.Signup(c.Request().Context(), User); err != nil {
		return
	}

	return c.JSON(http.StatusCreated, TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

func (eh *UserHandler) Login(c echo.Context) (err error) {
	User := &User{}
	if err = c.Bind(User); err != nil {
		return
	}

	var accessToken, refreshToken string
	if accessToken, refreshToken, err = eh.es.Login(c.Request().Context(), User.Login, User.Password); err != nil {
		return
	}

	return c.JSON(http.StatusOK, TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

func (eh *UserHandler) Refresh(c echo.Context) (err error) {
	token, login := tokenFromContext(c)

	var accessToken, refreshToken string
	if accessToken, refreshToken, err = eh.es.Refresh(c.Request().Context(), login, token); err != nil {
		return
	}

	return c.JSON(http.StatusOK, TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

func (eh *UserHandler) Update(c echo.Context) (err error) {
	User := &User{}
	if err = c.Bind(User); err != nil {
		return
	}
	_, login := tokenFromContext(c)

	if err = c.Validate(User); err != nil {
		return
	}

	if err = eh.es.Update(c.Request().Context(), login, User); err != nil {
		return
	}

	return c.JSON(http.StatusOK, login)
}

func (eh *UserHandler) Delete(c echo.Context) (err error) {
	_, login := tokenFromContext(c)

	if err = eh.es.Delete(c.Request().Context(), login); err != nil {
		return
	}

	return c.JSON(http.StatusOK, login)
}

func tokenFromContext(c echo.Context) (tokenRaw string, login string) {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	return user.Raw, claims["login"].(string)
}
