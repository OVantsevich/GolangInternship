// Package handler package with handlers
package handler

import (
	"context"
	"fmt"
	"net/http"

	"GolangInternship/FMicroservice/internal/model"
	"GolangInternship/FMicroservice/internal/service"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// UserClassicService service interface for user handler
//
//go:generate mockery --name=UserClassicService --case=underscore --output=./mocks
type UserClassicService interface {
	Signup(ctx context.Context, user *model.User) (string, string, *model.User, error)
	Login(ctx context.Context, login, password string) (string, string, error)
	Refresh(ctx context.Context, login, userRefreshToken string) (string, string, error)
	Update(ctx context.Context, login string, user *model.User) error
	Delete(ctx context.Context, login string) error

	GetByLogin(ctx context.Context, login string) (*model.User, error)
}

// UserClassic handler
type UserClassic struct {
	s UserClassicService
}

type tokenResponse struct {
	AccessToken  string `json:"access" example:"eyJhbGciOiJIUzI1NiIsInR5cC6IkpXVCJ9.eyJsb2dpbiI6InRc3QxIiwiZXhwIjoxNjc1MDgwNjE3fQ.OIt5MGzpbo1vZT5aNRvPwZCpU_tx-lisT2W2eyh78"`
	RefreshToken string `json:"refresh" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJsb2dpiI6InRlc3QxIiwiZXhwIjoxNjc1MTE1NzE3fQ.UJ0HF6D4Hb7cLdDfQxg3Byzvb8hWEXwK2RaNWDH54"`
}

type signupResponse struct {
	*model.User
	*tokenResponse
}

// NewUserHandlerClassic new user handler
func NewUserHandlerClassic(s UserClassicService) *UserClassic {
	return &UserClassic{s: s}
}

// Signup godoc
//
// @Summary      Add new user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        body	body     model.User  true  "New user object"
// @Success      201	{object}	signupResponse
// @Failure      400
// @Failure      500
// @Router       /signup [post]
func (h *UserClassic) Signup(c echo.Context) (err error) {
	user := &model.User{}
	err = c.Bind(user)
	if err != nil {
		logrus.Error(fmt.Errorf("userHandler - Signup - Bind: %w", err))
		return err
	}

	err = c.Validate(user)
	if err != nil {
		logrus.Error(fmt.Errorf("userHandler - Signup - Validate: %w", err))
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	var accessToken, refreshToken string
	var user2 *model.User
	accessToken, refreshToken, user2, err = h.s.Signup(c.Request().Context(), user)
	if err != nil {
		logrus.Error(fmt.Errorf("userHandler - Signup - Signup: %w", err))
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return c.JSON(http.StatusCreated,
		signupResponse{
			user2,
			&tokenResponse{
				AccessToken:  accessToken,
				RefreshToken: refreshToken,
			}})
}

// Login godoc
//
// @Summary      Login user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        body		body    model.User	true  "login and password"
// @Success      201	{object}	tokenResponse
// @Failure      500
// @Router       /login [post]
func (h *UserClassic) Login(c echo.Context) (err error) {
	user := &model.User{}
	err = c.Bind(user)
	if err != nil {
		logrus.Error(fmt.Errorf("userHandler - Login - Bind: %w", err))
		return err
	}

	var accessToken, refreshToken string
	if accessToken, refreshToken, err = h.s.Login(c.Request().Context(), user.Login, user.Password); err != nil {
		logrus.Error(fmt.Errorf("userHandler - Login - Login: %w", err))
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return c.JSON(http.StatusOK, tokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

// Refresh godoc
//
// @Summary      Refresh accessToken and refreshToken
// @Tags         users
// @Produce      json
// @Success      201	{object}	tokenResponse
// @Failure      500
// @Router       /refresh [get]
// @Security Bearer
func (h *UserClassic) Refresh(c echo.Context) (err error) {
	token, login := tokenFromContext(c)

	var accessToken, refreshToken string
	accessToken, refreshToken, err = h.s.Refresh(c.Request().Context(), login, token)
	if err != nil {
		logrus.Error(fmt.Errorf("userHandler - Refresh - Refresh: %w", err))
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return c.JSON(http.StatusOK, tokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

// Update godoc
//
// @Summary      Update info about user
// @Tags         users
// @Produce      json
// @Param		 body	body	model.User	 true	"New data"
// @Success      201	{string} string "login"
// @Failure      400
// @Failure      500
// @Router       /update [put]
// @Security Bearer
func (h *UserClassic) Update(c echo.Context) (err error) {
	user := &model.User{}
	err = c.Bind(user)
	if err != nil {
		logrus.Error(fmt.Errorf("userHandler - Update - Bind: %w", err))
		return err
	}
	_, login := tokenFromContext(c)

	err = c.Validate(user)
	if err != nil {
		logrus.Error(fmt.Errorf("userHandler - Update - Validate: %w", err))
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	err = h.s.Update(c.Request().Context(), login, user)
	if err != nil {
		logrus.Error(fmt.Errorf("userHandler - Update - Update: %w", err))
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return c.JSON(http.StatusOK, login)
}

// Delete godoc
//
// @Summary      Delete user
// @Tags         users
// @Produce      json
// @Success      201	{string} string "login"
// @Failure      500
// @Router       /delete [delete]
// @Security Bearer
func (h *UserClassic) Delete(c echo.Context) (err error) {
	_, login := tokenFromContext(c)

	err = h.s.Delete(c.Request().Context(), login)
	if err != nil {
		logrus.Error(fmt.Errorf("userHandler - Delete - Delete: %w", err))
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return c.JSON(http.StatusOK, login)
}

// UserByLogin godoc
//
// @Summary		 getting user by login
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        login	 header   string	true  "login"
// @Success      201	object	GetByLogin
// @Failure      403
// @Failure      500
// @Router       /admin/userByLogin [get]
// @Security Bearer
func (h *UserClassic) UserByLogin(c echo.Context) (err error) {
	login := c.Request().Header.Get("login")

	var user *model.User
	user, err = h.s.GetByLogin(c.Request().Context(), login)
	if err != nil {
		logrus.Error(fmt.Errorf("userHandler - UserByLogin - GetLastN: %w", err))
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return c.JSON(http.StatusOK, user)
}

func tokenFromContext(c echo.Context) (tokenRaw, login string) {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims
	return user.Raw, claims.(*service.CustomClaims).Login
}
