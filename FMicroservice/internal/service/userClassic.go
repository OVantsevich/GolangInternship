// Package service package with services
package service

import (
	"context"
	"fmt"
	"time"

	"GolangInternship/FMicroservice/internal/model"

	"github.com/golang-jwt/jwt/v4"
	passwordvalidator "github.com/wagslane/go-password-validator"
)

// UserClassicStream stream interface for user service
//
//go:generate mockery --name=UserClassicStream --case=underscore --output=./mocks
type UserClassicStream interface {
	ProduceUser(ctx context.Context, user *model.User) error
}

// UserClassicCache cache interface for user service
//
//go:generate mockery --name=UserClassicCache --case=underscore --output=./mocks
type UserClassicCache interface {
	GetByLogin(ctx context.Context, login string) (*model.User, bool, error)
	CreateUser(ctx context.Context, user *model.User) error
}

// UserClassicRepository repository interface for user service
//
//go:generate mockery --name=UserClassicRepository --case=underscore --output=./mocks
type UserClassicRepository interface {
	CreateUser(ctx context.Context, user *model.User) (*model.User, error)
	GetUserByLogin(ctx context.Context, login string) (*model.User, error)
	UpdateUser(ctx context.Context, login string, user *model.User) error
	RefreshUser(ctx context.Context, login, token string) error
	DeleteUser(ctx context.Context, login string) error
}

// Expiration time of access token
const accessExp = time.Minute * 15

// Expiration time of refresh token
const refreshExp = time.Hour * 10

// Strength of password
const passwordStrength = 50

// UserClassic user service
type UserClassic struct {
	rps    UserClassicRepository
	cache  UserClassicCache
	stream UserClassicStream
	jwtKey []byte
}

// CustomClaims claims with login and role
type CustomClaims struct {
	Login string `json:"login"`
	Role  string `json:"role"`
	jwt.RegisteredClaims
}

// NewUserServiceClassic new user service
func NewUserServiceClassic(rps UserClassicRepository, cache UserClassicCache, stream UserClassicStream, key string) *UserClassic {
	return &UserClassic{rps: rps, cache: cache, stream: stream, jwtKey: []byte(key)}
}

// Signup service signup
func (u *UserClassic) Signup(ctx context.Context, user *model.User) (accessToken, refreshToken string, user2 *model.User, err error) {
	if err = passwordvalidator.Validate(user.Password, passwordStrength); err != nil {
		return "", "", nil, fmt.Errorf("userService - Signup - Validate: %w", err)
	}

	if user2, err = u.rps.CreateUser(ctx, user); err != nil {
		return "", "", nil, fmt.Errorf("userService - Signup - CreateUser: %w", err)
	}

	user.Role = "user"
	accessToken, refreshToken, err = u.createJWT(ctx, user)
	if err != nil {
		return "", "", nil, fmt.Errorf("userService - Signup - createJWT: %w", err)
	}

	err = u.stream.ProduceUser(ctx, user)
	if err != nil {
		return "", "", nil, fmt.Errorf("userService - Signup - ProduceUser: %w", err)
	}

	return
}

// Login service login
//
//nolint:dupl //just because
func (u *UserClassic) Login(ctx context.Context, login, password string) (accessToken, refreshToken string, err error) {
	var user *model.User

	if user, err = u.rps.GetUserByLogin(ctx, login); err != nil {
		return "", "", fmt.Errorf("userService - Login - GetUserByLogin: %w", err)
	}

	if user.Password != password {
		return "", "", fmt.Errorf("userService - Login - Password invalid: %w", err)
	}

	accessToken, refreshToken, err = u.createJWT(ctx, user)
	if err != nil {
		return "", "", fmt.Errorf("userService - Login - createJWT: %w", err)
	}

	return
}

// Refresh service refresh
//
//nolint:dupl //just because
func (u *UserClassic) Refresh(ctx context.Context, login, userRefreshToken string) (accessToken, refreshToken string, err error) {
	var user *model.User

	if user, err = u.rps.GetUserByLogin(ctx, login); err != nil {
		return "", "", fmt.Errorf("userService - Refresh - GetUserByLogin: %w", err)
	}

	if user.Token != userRefreshToken {
		return "", "", fmt.Errorf("userService - Refresh - Token invalid: %w", err)
	}

	accessToken, refreshToken, err = u.createJWT(ctx, user)
	if err != nil {
		return "", "", fmt.Errorf("userService - Refresh - createJWT: %w", err)
	}

	return
}

// Update service update
func (u *UserClassic) Update(ctx context.Context, login string, user *model.User) (err error) {
	if err = u.rps.UpdateUser(ctx, login, user); err != nil {
		return fmt.Errorf("userService - Update - UpdateUser: %w", err)
	}

	return
}

// Delete service delete
func (u *UserClassic) Delete(ctx context.Context, login string) (err error) {
	if err = u.rps.DeleteUser(ctx, login); err != nil {
		return fmt.Errorf("userService - Delete - DeleteUser: %w", err)
	}

	return
}

// GetByLogin service get by login
func (u *UserClassic) GetByLogin(ctx context.Context, login string) (user *model.User, err error) {
	var notCached bool
	user, notCached, err = u.cache.GetByLogin(ctx, login)
	if err != nil && !notCached {
		return nil, fmt.Errorf("userService - GetByLogin - cache - GetByLogin: %w", err)
	}

	if notCached {
		if user, err = u.rps.GetUserByLogin(ctx, login); err != nil {
			return nil, fmt.Errorf("userService - GetByLogin - Repository - GetByLogin: %w", err)
		}
		err = u.cache.CreateUser(ctx, user)
		if err != nil {
			return nil, fmt.Errorf("userService - GetByLogin - cache - CreateUser: %w", err)
		}
		return
	}

	return
}

func (u *UserClassic) createJWT(ctx context.Context, user *model.User) (accessTokenStr, refreshTokenStr string, err error) {
	accessClaims := &CustomClaims{
		user.Login,
		user.Role,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessExp)),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenStr, err = accessToken.SignedString(u.jwtKey)
	if err != nil {
		return "", "", fmt.Errorf("userService - createJWT - SignedString: %w", err)
	}

	refreshClaims := &CustomClaims{
		user.Login,
		user.Role,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshExp)),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenStr, err = refreshToken.SignedString(u.jwtKey)
	if err != nil {
		return "", "", fmt.Errorf("userService - createJWT - SignedString: %w", err)
	}

	err = u.rps.RefreshUser(ctx, user.Login, refreshTokenStr)
	if err != nil {
		return "", "", fmt.Errorf("userService - createJWT - RefreshUser: %w", err)
	}
	return accessTokenStr, refreshTokenStr, err
}
