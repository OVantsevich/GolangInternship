package service

import (
	"context"
	"github.com/OVantsevich/GolangInternship/FMicroservice/internal/model"
	mocks "github.com/OVantsevich/GolangInternship/FMicroservice/mocks/repository"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

var service *User

var testValidData = []model.User{
	{
		Name:     `NAME`,
		Age:      1,
		Login:    `CreateLOGIN1`,
		Email:    `LOGIN1@gmail.com`,
		Password: `strongPassword`,
	},
	{
		Name:     `NAME`,
		Age:      1,
		Login:    `CreateLOGIN2`,
		Email:    `LOGIN2@gmail.com`,
		Password: `PASSWORD123456789`,
	},
}
var testNoValidData = []model.User{
	{
		Name:     `nameEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEE`,
		Age:      22222,
		Login:    `LOGIN2`,
		Email:    `LOGIN2@gmail.com`,
		Password: `PASSWORD123`,
	},
	{
		Name:     `NAME`,
		Age:      2,
		Login:    `LOGIN1AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA`,
		Email:    `LOGIN1@gmail.com`,
		Password: `LOGIN23102002`,
	},
}

func TestUser_Signup(t *testing.T) {
	repository := mocks.NewUser(t)
	service = NewUserService(repository, "secret-key")
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

}
