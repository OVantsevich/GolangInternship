package handler

import (
	"fmt"
	"github.com/OVantsevich/GolangInternship/FMicroservice/internal/model"
	"github.com/OVantsevich/GolangInternship/FMicroservice/internal/service"
	"github.com/golang-jwt/jwt/v4"
	echo "github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
)

type UserClassic struct {
	s service.User
}

type TokenResponse struct {
	AccessToken  string `json:"access" example:"eyJhbGciOiJIUzI1NiIsInR5cC6IkpXVCJ9.eyJsb2dpbiI6InRc3QxIiwiZXhwIjoxNjc1MDgwNjE3fQ.OIt5MGzpbo1vZT5aNRvPwZCpU_tx-lisT2W2eyh78"`
	RefreshToken string `json:"refresh" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJsb2dpiI6InRlc3QxIiwiZXhwIjoxNjc1MTE1NzE3fQ.UJ0HF6D4Hb7cLdDfQxg3Byzvb8hWEXwK2RaNWDH54"`
}

type SignupResponse struct {
	*model.User
	*TokenResponse
}

func NewUserHandlerClassic(s service.User) *UserClassic {
	return &UserClassic{s: s}
}

// Signup godoc
//
// @Summary      Add new user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        body	body     model.User  true  "New user object"
// @Success      201	{object}	SignupResponse
// @Failure      400
// @Failure      500
// @Router       /signup [post]
func (h *UserClassic) Signup(c echo.Context) (err error) {
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

// Login godoc
//
// @Summary      Login user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        body		body    model.User	true  "login and password"
// @Success      201	{object}	TokenResponse
// @Failure      500
// @Router       /login [post]
func (h *UserClassic) Login(c echo.Context) (err error) {
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

// Refresh godoc
//
// @Summary      Refresh accessToken and refreshToken
// @Tags         users
// @Produce      json
// @Success      201	{object}	TokenResponse
// @Failure      500
// @Router       /refresh [get]
// @Security Bearer
func (h *UserClassic) Refresh(c echo.Context) (err error) {
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
