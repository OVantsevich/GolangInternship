package handler

import (
	"GolangInternship/FMicroserviceGRPC/internal/model"
	"GolangInternship/FMicroserviceGRPC/internal/service"
	pr "GolangInternship/FMicroserviceGRPC/proto"
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
)

//go:generate mockery --name=UserClassicService --case=underscore --output=./mocks
type UserClassicService interface {
	Signup(ctx context.Context, user *model.User) (string, string, *model.User, error)
	Login(ctx context.Context, login, password string) (string, string, error)
	Refresh(ctx context.Context, login, userRefreshToken string) (string, string, error)
	Update(ctx context.Context, login string, user *model.User) error
	Delete(ctx context.Context, login string) error

	GetByLogin(ctx context.Context, login string) (*model.User, error)
}

type UserClassic struct {
	pr.UnimplementedUserServiceServer
	s      UserClassicService
	jwtKey string
}

type TokenResponse struct {
	AccessToken  string `json:"access" example:"eyJhbGciOiJIUzI1NiIsInR5cC6IkpXVCJ9.eyJsb2dpbiI6InRc3QxIiwiZXhwIjoxNjc1MDgwNjE3fQ.OIt5MGzpbo1vZT5aNRvPwZCpU_tx-lisT2W2eyh78"`
	RefreshToken string `json:"refresh" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJsb2dpiI6InRlc3QxIiwiZXhwIjoxNjc1MTE1NzE3fQ.UJ0HF6D4Hb7cLdDfQxg3Byzvb8hWEXwK2RaNWDH54"`
}

type SignupResponse struct {
	*model.User
	*TokenResponse
}

func NewUserHandlerClassic(s UserClassicService, key string) *UserClassic {
	return &UserClassic{s: s, jwtKey: key}
}

func (h *UserClassic) Signup(ctx context.Context, request *pr.SignupRequest) (response *pr.SignupResponse, err error) {
	user := &model.User{
		Login:    request.Login,
		Email:    request.Email,
		Password: request.Password,
		Name:     request.Name,
		Age:      int(request.Age),
	}

	var userResponse *model.User
	response = &pr.SignupResponse{}
	if response.AccessToken, response.RefreshToken, userResponse, err = h.s.Signup(ctx, user); err != nil {
		err = fmt.Errorf("userHandler - Signup - Signup: %w", err)
		logrus.Error(err)
		return
	}
	response.User = &pr.User{
		Login:    userResponse.Login,
		Email:    userResponse.Email,
		Password: userResponse.Password,
		Name:     userResponse.Name,
		Age:      int32(userResponse.Age),
		Role:     userResponse.Role,
	}

	return
}

func (h *UserClassic) Login(ctx context.Context, request *pr.LoginRequest) (response *pr.LoginResponse, err error) {
	response = &pr.LoginResponse{}
	if response.AccessToken, response.RefreshToken, err = h.s.Login(ctx, request.Login, request.Password); err != nil {
		err = fmt.Errorf("userHandler - Login - Login: %w", err)
		logrus.Error(err)
		return
	}

	return
}

func (h *UserClassic) Refresh(ctx context.Context, request *pr.RefreshRequest) (response *pr.RefreshResponse, err error) {
	response = &pr.RefreshResponse{}
	if response.AccessToken, response.RefreshToken, err = h.s.Refresh(ctx, request.Login, request.RefreshToken); err != nil {
		err = fmt.Errorf("userHandler - Refresh - Refresh: %w", err)
		logrus.Error(err)
		return
	}

	return
}

func (h *UserClassic) Update(ctx context.Context, request *pr.UpdateRequest) (response *pr.UpdateResponse, err error) {
	var claims *service.CustomClaims
	claims, err = h.Verify(request.AccessToken)
	if err != nil {
		err = fmt.Errorf("userHandler - Update - Verify: %w", err)
		logrus.Error(err)
		return
	}

	user := &model.User{
		Login:    request.User.Login,
		Email:    request.User.Email,
		Password: request.User.Password,
		Name:     request.User.Name,
		Age:      int(request.User.Age),
	}
	response = &pr.UpdateResponse{}
	if err = h.s.Update(ctx, claims.Login, user); err != nil {
		err = fmt.Errorf("userHandler - Update - Update: %w", err)
		logrus.Error(err)
		return
	}
	response.Login = claims.Login

	return
}

func (h *UserClassic) Delete(ctx context.Context, request *pr.DeleteRequest) (response *pr.DeleteResponse, err error) {
	var claims *service.CustomClaims
	claims, err = h.Verify(request.AccessToken)
	if err != nil {
		err = fmt.Errorf("userHandler - Delete - Verify: %w", err)
		logrus.Error(err)
		return
	}

	response = &pr.DeleteResponse{}
	if err = h.s.Delete(ctx, claims.Login); err != nil {
		err = fmt.Errorf("userHandler - Delete - Delete: %w", err)
		logrus.Error(err)
		return
	}
	response.Login = claims.Login

	return
}

func (h *UserClassic) UserByLogin(ctx context.Context, request *pr.UserByLoginRequest) (response *pr.UserByLoginResponse, err error) {
	var claims *service.CustomClaims
	claims, err = h.Verify(request.AccessToken)
	if err != nil {
		err = fmt.Errorf("userHandler - UserByLogin - Verify: %w", err)
		logrus.Error(err)
		return
	}
	if claims.Role != "admin" {
		err = fmt.Errorf("access denied")
		logrus.Error(err)
		return
	}

	response = &pr.UserByLoginResponse{}
	var user *model.User
	if user, err = h.s.GetByLogin(ctx, claims.Login); err != nil {
		err = fmt.Errorf("userHandler - UserByLogin - GetByLogin: %w", err)
		logrus.Error(err)
		return
	}
	response.User = &pr.User{
		Login:    user.Login,
		Email:    user.Email,
		Password: user.Password,
		Name:     user.Name,
		Age:      int32(user.Age),
		Role:     user.Role,
	}

	return
}

func (h *UserClassic) Verify(token string) (claims *service.CustomClaims, err error) {
	claims = &service.CustomClaims{}

	_, err = jwt.ParseWithClaims(
		token,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(h.jwtKey), nil
		},
	)
	if err != nil {
		err = fmt.Errorf("invalid token: %w", err)
	}

	return
}
