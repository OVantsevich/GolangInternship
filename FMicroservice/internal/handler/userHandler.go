package handler

import (
	"fmt"
	"github.com/OVantsevich/GolangInternship/FMicroservice/internal/model"
	"github.com/OVantsevich/GolangInternship/FMicroservice/internal/service"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
)

type User struct {
	s *service.User
}

type TokenResponse struct {
	AccessToken  string `json:"access"`
	RefreshToken string `json:"refresh"`
}

type SignupResponse struct {
	*model.User
	*TokenResponse
}

func NewUserHandler(s *service.User) *User {
	return &User{s: s}
}

func (h *User) Signup(c echo.Context) (err error) {
	user := &model.User{}
	if err = c.Bind(user); err != nil {
		logrus.Error(fmt.Errorf("userHandler - Signup - Bind: %w", err))
		return
	}

	if err = c.Validate(user); err != nil {
		logrus.Error(fmt.Errorf("userHandler - Signup - Validate: %w", err))
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	var accessToken, refreshToken string
	var user2 *model.User
	if accessToken, refreshToken, user2, err = h.s.Signup(c.Request().Context(), user); err != nil {
		logrus.Error(fmt.Errorf("userHandler - Signup - Signup: %w", err))
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return c.JSON(http.StatusCreated,
		SignupResponse{
			user2,
			&TokenResponse{
				AccessToken:  accessToken,
				RefreshToken: refreshToken,
			}})
}

func (h *User) Login(c echo.Context) (err error) {
	user := &model.User{}
	if err = c.Bind(user); err != nil {
		logrus.Error(fmt.Errorf("userHandler - Login - Bind: %w", err))
		return
	}

	var accessToken, refreshToken string
	if accessToken, refreshToken, err = h.s.Login(c.Request().Context(), user.Login, user.Password); err != nil {
		logrus.Error(fmt.Errorf("userHandler - Login - Login: %w", err))
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return c.JSON(http.StatusOK, TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

func (h *User) Refresh(c echo.Context) (err error) {
	token, login := tokenFromContext(c)

	var accessToken, refreshToken string
	if accessToken, refreshToken, err = h.s.Refresh(c.Request().Context(), login, token); err != nil {
		logrus.Error(fmt.Errorf("userHandler - Refresh - Refresh: %w", err))
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return c.JSON(http.StatusOK, TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

func (h *User) Update(c echo.Context) (err error) {
	user := &model.User{}
	if err = c.Bind(user); err != nil {
		logrus.Error(fmt.Errorf("userHandler - Update - Bind: %w", err))
		return
	}
	_, login := tokenFromContext(c)

	if err = c.Validate(user); err != nil {
		logrus.Error(fmt.Errorf("userHandler - Update - Validate: %w", err))
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	if err = h.s.Update(c.Request().Context(), login, user); err != nil {
		logrus.Error(fmt.Errorf("userHandler - Update - Update: %w", err))
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return c.JSON(http.StatusOK, login)
}

func (h *User) Delete(c echo.Context) (err error) {
	_, login := tokenFromContext(c)

	if err = h.s.Delete(c.Request().Context(), login); err != nil {
		logrus.Error(fmt.Errorf("userHandler - Delete - Delete: %w", err))
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return c.JSON(http.StatusOK, login)
}

func tokenFromContext(c echo.Context) (tokenRaw string, login string) {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims
	return user.Raw, claims.(*service.CustomClaims).Login
}
