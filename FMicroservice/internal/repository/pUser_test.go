package repository

import (
	"context"
	"github.com/OVantsevich/GolangInternship/FMicroservice/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"testing"
)

var prps *PUser

var testValidData = []model.User{
	{
		Name:     `NAME`,
		Age:      1,
		Login:    `CreateLOGIN11`,
		Email:    `LOGIN1@gmail.com`,
		Password: `LOGIN123456789`,
	},
	{
		Name:     `NAME`,
		Age:      1,
		Login:    `CreateLOGIN22`,
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

func NewPRepository(pool *pgxpool.Pool) *PUser {
	return &PUser{Pool: pool}
}

func TestPUser_CreateUser(t *testing.T) {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, "postgres://postgres:postgres@localhost:5432/userService?sslmode=disable")
	require.NoError(t, err, "new pool error")
	prps = NewPRepository(pool)

	for _, u := range testValidData {
		_, err = prps.CreateUser(ctx, &u)
		require.NoError(t, err, "create error")

		_, err = prps.Pool.Exec(ctx, "delete from users where login=$1 ", u.Login)
	}

	for _, u := range testNoValidData {
		_, err = prps.CreateUser(ctx, &u)
		require.Error(t, err, "create error")
	}

	// Already existing data
	for _, u := range testValidData {
		_, err = prps.CreateUser(ctx, &u)
		require.NoError(t, err, "create error")

		_, err = prps.CreateUser(ctx, &u)
		require.Error(t, err, "create error")

		_, err = prps.Pool.Exec(ctx, "delete from users where login=$1 ", u.Login)
	}
}

func TestPUser_GetUserByLogin(t *testing.T) {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, "postgres://postgres:postgres@localhost:5432/userService?sslmode=disable")
	require.NoError(t, err, "new pool error")
	prps := NewPRepository(pool)

	var user *model.User
	for _, u := range testValidData {
		_, err = prps.Pool.Exec(ctx, "delete from users where login=$1 ", u.Login)
		_, err = prps.CreateUser(ctx, &u)
		require.NoError(t, err, "create error")

		user, _, err = prps.GetUserByLogin(ctx, u.Login)
		require.Equal(t, u.Password, user.Password)
		require.Equal(t, u.Email, user.Email)
		require.NoError(t, err, "get by login error")

		_, err = prps.Pool.Exec(ctx, "delete from users where login=$1 ", u.Login)
	}

	//Non-existent data
	for _, u := range testValidData {
		_, err = prps.Pool.Exec(ctx, "delete from users where login=$1 ", u.Login)

		user, _, err = prps.GetUserByLogin(ctx, u.Login)
		require.Error(t, err, "get by login error")
	}
}

func TestPUser_UpdateUser(t *testing.T) {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, "postgres://postgres:postgres@localhost:5432/userService?sslmode=disable")
	require.NoError(t, err, "new pool error")
	prps := NewPRepository(pool)

	var user *model.User
	for _, u := range testValidData {
		_, err = prps.Pool.Exec(ctx, "delete from users where login=$1 ", u.Login)
		_, err = prps.CreateUser(ctx, &u)
		require.NoError(t, err, "create error")

		u.Name = "Update"
		err = prps.UpdateUser(ctx, u.Login, &u)
		require.NoError(t, err, "update error")

		user, _, err = prps.GetUserByLogin(ctx, u.Login)
		require.Equal(t, "Update", user.Name)
		require.NoError(t, err, "get by login error")

		_, err = prps.Pool.Exec(ctx, "delete from users where login=$1 ", u.Login)
	}

	//Invalid data
	_, err = prps.Pool.Exec(ctx, "delete from users where login=$1 ", &testValidData[0].Login)
	_, err = prps.CreateUser(ctx, &testValidData[0])
	require.NoError(t, err, "create error")
	err = prps.UpdateUser(ctx, testValidData[0].Login, &testNoValidData[0])
	require.Error(t, err, "update error")
	_, err = prps.Pool.Exec(ctx, "delete from users where login=$1 ", &testValidData[0].Login)

	//Non-existent data
	for _, u := range testValidData {
		_, err = prps.Pool.Exec(ctx, "delete from users where login=$1 ", u.Login)

		err = prps.UpdateUser(ctx, u.Login, &u)
		require.Error(t, err, "update error")
	}
}

func TestPUser_RefreshUser(t *testing.T) {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, "postgres://postgres:postgres@localhost:5432/userService?sslmode=disable")
	require.NoError(t, err, "new pool error")
	prps := NewPRepository(pool)

	var user *model.User
	for _, u := range testValidData {
		_, err = prps.Pool.Exec(ctx, "delete from users where login=$1 ", u.Login)
		_, err = prps.CreateUser(ctx, &u)
		require.NoError(t, err, "create error")

		token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJsb2dpbiI6InRlc3QxIiwiZXhwIjoxNjc0ODMxODE2fQ.jlD1_wrfdK8XjMut236sQDb7B7EOvVjflGZnNUS5o2g"
		err = prps.RefreshUser(ctx, u.Login, token)
		require.NoError(t, err, "refresh error")

		user, _, err = prps.GetUserByLogin(ctx, u.Login)
		require.Equal(t, token, user.Token)
		require.NoError(t, err, "get by login error")

		_, err = prps.Pool.Exec(ctx, "delete from users where login=$1 ", u.Login)
	}

	//Non-existent data
	for _, u := range testValidData {
		_, err = prps.Pool.Exec(ctx, "delete from users where login=$1 ", u.Login)

		token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJsb2dpbiI6InRlc3QxIiwiZXhwIjoxNjc0ODMxODE2fQ.jlD1_wrfdK8XjMut236sQDb7B7EOvVjflGZnNUS5o2g"
		err = prps.RefreshUser(ctx, u.Login, token)
		require.Error(t, err, "refresh error")
	}
}

func TestPUser_DeleteUser(t *testing.T) {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, "postgres://postgres:postgres@localhost:5432/userService?sslmode=disable")
	require.NoError(t, err, "new pool error")
	prps := NewPRepository(pool)

	for _, u := range testValidData {
		_, err = prps.Pool.Exec(ctx, "delete from users where login=$1 ", u.Login)
		_, err = prps.CreateUser(ctx, &u)
		require.NoError(t, err, "create error")

		err = prps.DeleteUser(ctx, u.Login)
		require.NoError(t, err, "delete error")

		_, _, err = prps.GetUserByLogin(ctx, u.Login)
		require.Error(t, err, "get by login error")

		_, err = prps.Pool.Exec(ctx, "delete from users where login=$1 ", u.Login)
	}

	//Non-existent data
	for _, u := range testValidData {
		_, err = prps.Pool.Exec(ctx, "delete from users where login=$1 ", u.Login)

		err = prps.DeleteUser(ctx, u.Login)
		require.Error(t, err, "delete error")
	}
}
