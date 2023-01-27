package service

import (
	"context"
	"fmt"
	"github.com/OVantsevich/GolangInternship/FMicroservice/internal/model"
	"github.com/OVantsevich/GolangInternship/FMicroservice/internal/repository"
	"github.com/golang-jwt/jwt/v4"
	passwordvalidator "github.com/wagslane/go-password-validator"
	"time"
)

type User struct {
	rps    repository.User
	jwtKey []byte
}

type CustomClaims struct {
	Login string `json:"login"`
	jwt.RegisteredClaims
}

func NewUserService(rps repository.User, key string) *User {
	return &User{rps: rps, jwtKey: []byte(key)}
}

func (u *User) Signup(ctx context.Context, user *model.User) (accessToken, refreshToken string, user2 *model.User, err error) {

	if err = passwordvalidator.Validate(user.Password, 50); err != nil {
		return "", "", nil, fmt.Errorf("userService - Signup - Validate: %w", err)
	}

	if user2, err = u.rps.CreateUser(ctx, user); err != nil {
		return "", "", nil, fmt.Errorf("userService - Signup - CreateUser: %w", err)
	}

	accessToken, refreshToken, err = u.CreateJWT(ctx, user)
	if err != nil {
		return "", "", nil, fmt.Errorf("userService - Signup - CreateJWT: %w", err)
	}

	return
}

func (u *User) Login(ctx context.Context, login, password string) (accessToken, refreshToken string, err error) {
	var user *model.User
	if user, err = u.rps.GetUserByLogin(ctx, login); err != nil {
		return "", "", fmt.Errorf("userService - Login - GetUserByLogin: %w", err)
	}

	if user.Password != password {
		return "", "", fmt.Errorf("userService - Login - Password invalid: %w", err)
	}

	accessToken, refreshToken, err = u.CreateJWT(ctx, user)
	if err != nil {
		return "", "", fmt.Errorf("userService - Login - CreateJWT: %w", err)
	}

	return
}

func (u *User) Refresh(ctx context.Context, login, userRefreshToken string) (accessToken, refreshToken string, err error) {
	var user *model.User

	if user, err = u.rps.GetUserByLogin(ctx, login); err != nil {
		return "", "", fmt.Errorf("userService - Refresh - GetUserByLogin: %w", err)
	}

	if user.Token != userRefreshToken {
		return "", "", fmt.Errorf("userService - Refresh - Token invalid: %w", err)
	}

	accessToken, refreshToken, err = u.CreateJWT(ctx, user)
	if err != nil {
		return "", "", fmt.Errorf("userService - Refresh - CreateJWT: %w", err)
	}

	return
}

func (u *User) Update(ctx context.Context, login string, user *model.User) (err error) {
	if err = u.rps.UpdateUser(ctx, login, user); err != nil {
		return fmt.Errorf("userService - Update - UpdateUser: %w", err)
	}

	return
}

func (u *User) Delete(ctx context.Context, login string) (err error) {
	if err = u.rps.DeleteUser(ctx, login); err != nil {
		return fmt.Errorf("userService - Delete - DeleteUser: %w", err)
	}

	return
}

func (u *User) CreateJWT(ctx context.Context, user *model.User) (accessTokenStr, refreshTokenStr string, err error) {
	accessClaims := &CustomClaims{
		user.Login,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenStr, err = accessToken.SignedString(u.jwtKey)
	if err != nil {
		return "", "", fmt.Errorf("userService - CreateJWT - SignedString: %w", err)
	}

	refreshClaims := &CustomClaims{
		user.Login,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 10)),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenStr, err = refreshToken.SignedString(u.jwtKey)
	if err != nil {
		return "", "", fmt.Errorf("userService - CreateJWT - SignedString: %w", err)
	}

	err = u.rps.RefreshUser(ctx, user.Login, refreshTokenStr)
	if err != nil {
		return "", "", fmt.Errorf("userService - CreateJWT - RefreshUser: %w", err)
	}
	return
}
