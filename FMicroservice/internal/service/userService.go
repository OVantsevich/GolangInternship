package service

import (
	"context"
	"fmt"
	. "github.com/OVantsevich/GolangInternship/FMicroservice/internal/domain"
	. "github.com/OVantsevich/GolangInternship/FMicroservice/internal/repository"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	passwordvalidator "github.com/wagslane/go-password-validator"
	"net/http"
	"time"
)

type UserService struct {
	rps    Repository
	jwtKey []byte
}

type CustomClaims struct {
	Login string
	jwt.RegisteredClaims
}

func NewUserService(rps *Repository) *UserService {
	return &UserService{rps: *rps}
}

func (us *UserService) Signup(ctx context.Context, user *User) (accessToken, refreshToken string, err error) {

	if err = passwordvalidator.Validate(user.Password, 50); err != nil {
		return "", "", &echo.HTTPError{Code: http.StatusBadRequest, Message: err.Error()}
	}

	if err = us.rps.CreateUser(ctx, user); err != nil {
		return "", "", &echo.HTTPError{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	accessToken, refreshToken, err = us.CreateJWT(ctx, user)
	if err != nil {
		return "", "", &echo.HTTPError{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	return
}

func (us *UserService) Login(ctx context.Context, login, password string) (accessToken, refreshToken string, err error) {
	var user *User
	if user, err = us.rps.GetUserByLogin(ctx, login); err != nil {
		return "", "", &echo.HTTPError{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	if user.Password != password {
		return "", "", &echo.HTTPError{Code: http.StatusBadRequest, Message: "invalid password"}
	}

	accessToken, refreshToken, err = us.CreateJWT(ctx, user)
	if err != nil {
		return "", "", &echo.HTTPError{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	return
}

func (us *UserService) Refresh(ctx context.Context, login, userRefreshToken string) (accessToken, refreshToken string, err error) {
	var user *User

	if user, err = us.rps.GetUserByLogin(ctx, login); err != nil {
		return "", "", &echo.HTTPError{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	if user.Token != userRefreshToken {
		return "", "", &echo.HTTPError{Code: http.StatusBadRequest, Message: "invalid token"}
	}

	accessToken, refreshToken, err = us.CreateJWT(ctx, user)
	if err != nil {
		return "", "", &echo.HTTPError{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	return
}

func (us *UserService) Update(ctx context.Context, login string, user *User) (err error) {
	if err = us.rps.UpdateUser(ctx, login, user); err != nil {
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	return
}

func (us *UserService) Delete(ctx context.Context, login string) (err error) {
	if err = us.rps.DeleteUser(ctx, login); err != nil {
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	return
}

type Claims struct {
	jwt.RegisteredClaims
}

func (us *UserService) CreateJWT(ctx context.Context, user *User) (accessTokenStr, refreshTokenStr string, err error) {
	accessClaims := &CustomClaims{
		user.Login,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenStr, err = accessToken.SignedString(us.jwtKey)
	if err != nil {
		return "", "", fmt.Errorf("service - userService - CreateJWT: %v", err)
	}

	refreshClaims := &CustomClaims{
		user.Login,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenStr, err = refreshToken.SignedString(us.jwtKey)
	if err != nil {
		return "", "", fmt.Errorf("service - userService - CreateJWT: %v", err)
	}

	err = us.rps.RefreshUser(ctx, user.Login, refreshTokenStr)
	if err != nil {
		return "", "", fmt.Errorf("service - userService - CreateJWT: %v", err)
	}
	return
}
