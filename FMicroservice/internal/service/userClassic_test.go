package service

import (
	"context"
	"fmt"
	"github.com/OVantsevich/GolangInternship/FMicroservice/internal/model"
	mocks "github.com/OVantsevich/GolangInternship/FMicroservice/mocks/repository"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

var service *UserClassic

var testValidData = []model.User{
	{
		Name:     `NAME`,
		Age:      1,
		Login:    `CreateLOGIN1`,
		Email:    `LOGIN1@gmail.com`,
		Token:    `validToken`,
		Password: `strongPassword`,
	},
	{
		Name:     `NAME`,
		Age:      1,
		Login:    `CreateLOGIN2`,
		Email:    `LOGIN2@gmail.com`,
		Token:    `validToken2`,
		Password: `PASSWORD123456789`,
	},
}
var testNoValidData = []model.User{
	{
		Name:     `nameEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEE`,
		Age:      22222,
		Login:    `LOGIN2`,
		Email:    `LOGIN2@gmail.com`,
		Token:    `dafrawerfaegfaegae`,
		Password: `weak`,
	},
	{
		Name:     `NAME`,
		Age:      2,
		Login:    `LOGIN1AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA`,
		Email:    `LOGIN1@gmail.com`,
		Token:    `argawegfafawfew`,
		Password: `lalala`,
	},
}

func TestUser_Signup(t *testing.T) {
	repository := mocks.NewUser(t)
	service = NewUserServiceClassic(repository, "secret-key")
	repository.On("CreateUser", mock.AnythingOfType("*context.emptyCtx"),
		mock.AnythingOfType("*model.User")).Return(nil, nil)
	repository.On("RefreshUser", mock.AnythingOfType("*context.emptyCtx"),
		mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)

	ctx := context.Background()
	var err error

	var rt, at string
	for _, user := range testValidData {
		at, rt, _, err = service.Signup(ctx, &user)
		require.NoError(t, err)

		claims := &CustomClaims{}
		_, err = jwt.ParseWithClaims(at, claims, func(token *jwt.Token) (interface{}, error) {
			return service.jwtKey, nil
		})
		require.NoError(t, err)
		require.NoError(t, claims.Valid())

		_, err = jwt.ParseWithClaims(rt, claims, func(token *jwt.Token) (interface{}, error) {
			return service.jwtKey, nil
		})
		require.NoError(t, err)
		require.NoError(t, claims.Valid())
	}

	for _, user := range testNoValidData {
		at, rt, _, err = service.Signup(ctx, &user)
		require.Error(t, err)
	}
}

func TestUser_Login(t *testing.T) {
	repository := mocks.NewUser(t)
	service = NewUserServiceClassic(repository, "secret-key")
	repository.On("RefreshUser", mock.AnythingOfType("*context.emptyCtx"),
		mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)

	ctx := context.Background()
	var err error

	for _, user := range testValidData {
		repository.On("GetUserByLogin", mock.AnythingOfType("*context.emptyCtx"),
			mock.AnythingOfType("string")).
			Return(func(ctx context.Context, s string) *model.User {
				return &model.User{Password: user.Password}
			},
				nil).Once()

		_, _, err = service.Login(ctx, user.Login, user.Password)
		require.NoError(t, err)
	}

	for _, user := range testNoValidData {
		repository.On("GetUserByLogin", mock.AnythingOfType("*context.emptyCtx"),
			mock.AnythingOfType("string")).
			Return(func(ctx context.Context, s string) *model.User {
				return &model.User{Password: "user.Password"}
			},
				nil).Once()
		_, _, err = service.Login(ctx, user.Login, user.Password)
		require.Error(t, err)
	}
}

func TestUser_Refresh(t *testing.T) {
	repository := mocks.NewUser(t)
	service = NewUserServiceClassic(repository, "secret-key")
	repository.On("RefreshUser", mock.AnythingOfType("*context.emptyCtx"),
		mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)

	ctx := context.Background()
	var err error

	for _, user := range testValidData {
		repository.On("GetUserByLogin", mock.AnythingOfType("*context.emptyCtx"),
			mock.AnythingOfType("string")).
			Return(func(ctx context.Context, s string) *model.User {
				return &model.User{Token: user.Token}
			},
				nil).Once()

		_, _, err = service.Refresh(ctx, user.Login, user.Token)
		require.NoError(t, err)
	}

	for _, user := range testNoValidData {
		repository.On("GetUserByLogin", mock.AnythingOfType("*context.emptyCtx"),
			mock.AnythingOfType("string")).
			Return(func(ctx context.Context, s string) *model.User {
				return &model.User{Token: "user.Token"}
			},
				nil).Once()
		_, _, err = service.Refresh(ctx, user.Login, user.Token)
		require.Error(t, err)
	}
}

func TestUser_Update(t *testing.T) {
	repository := mocks.NewUser(t)
	service = NewUserServiceClassic(repository, "secret-key")

	ctx := context.Background()
	var err error
	repository.On("UpdateUser", mock.AnythingOfType("*context.emptyCtx"),
		mock.AnythingOfType("string"), mock.AnythingOfType("*model.User")).Return(nil).Once()
	err = service.Update(ctx, testValidData[0].Login, &testValidData[0])
	require.NoError(t, err)

	repository.On("UpdateUser", mock.AnythingOfType("*context.emptyCtx"),
		mock.AnythingOfType("string"), mock.AnythingOfType("*model.User")).Return(fmt.Errorf("something went wrong")).Once()
	err = service.Update(ctx, testValidData[0].Login, &testValidData[0])
	require.Error(t, err)
}

func TestUser_Delete(t *testing.T) {
	repository := mocks.NewUser(t)
	service = NewUserServiceClassic(repository, "secret-key")

	ctx := context.Background()
	var err error
	repository.On("DeleteUser", mock.AnythingOfType("*context.emptyCtx"),
		mock.AnythingOfType("string")).Return(nil).Once()
	err = service.Delete(ctx, testValidData[0].Login)
	require.NoError(t, err)

	repository.On("DeleteUser", mock.AnythingOfType("*context.emptyCtx"),
		mock.AnythingOfType("string")).Return(fmt.Errorf("something went wrong")).Once()
	err = service.Delete(ctx, testValidData[0].Login)
	require.Error(t, err)
}
