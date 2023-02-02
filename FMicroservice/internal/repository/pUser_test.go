package repository

import (
	"context"
	"fmt"
	"github.com/OVantsevich/GolangInternship/FMicroservice/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ory/dockertest/v3"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"os"
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

var db *pgxpool.Pool

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		logrus.Fatalf("Could not construct pool: %s", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		logrus.Fatalf("Could not connect to Docker: %s", err)
	}

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "latest",
		Env: []string{
			"POSTGRES_USER=postgres",
			"POSTGRES_PASSWORD=postgres",
			"POSTGRES_DB=userService",
			"listen_addresses = '*'",
		},
		Mounts: []string{"/home/olegvantsevich/GolandProjects/GolangInternship/FMicroservice/migrations:/docker-entrypoint-initdb.d"},
	})
	if err != nil {
		logrus.Fatalf("Could not start resource: %s", err)
	}

	ctx := context.Background()
	if err := pool.Retry(func() error {
		var err error
		db, err = pgxpool.New(ctx, fmt.Sprintf("postgres://postgres:postgres@localhost:%s/userService?sslmode=disable", resource.GetPort("5432/tcp")))
		if err != nil {
			return err
		}
		return db.Ping(ctx)
	}); err != nil {
		logrus.Fatalf("Could not connect to database: %s", err)
	}
	db.Exec(ctx, "insert into roles (name) values ('admin')")
	db.Exec(ctx, "insert into roles (name) values ('user')")
	code := m.Run()

	if err := pool.Purge(resource); err != nil {
		logrus.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func NewPRepository(pool *pgxpool.Pool) *PUser {
	return &PUser{Pool: pool}
}

func TestPUser_CreateUser(t *testing.T) {
	ctx := context.Background()
	prps = NewPRepository(db)
	var err error
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
	prps := NewPRepository(db)

	var user *model.User
	var err error
	for _, u := range testValidData {
		_, err = prps.Pool.Exec(ctx, "delete from users where login=$1 ", u.Login)
		_, err = prps.CreateUser(ctx, &u)
		require.NoError(t, err, "create error")

		user, err = prps.GetUserByLogin(ctx, u.Login)
		require.Equal(t, u.Password, user.Password)
		require.Equal(t, u.Email, user.Email)
		require.NoError(t, err, "get by login error")

		_, err = prps.Pool.Exec(ctx, "delete from users where login=$1 ", u.Login)
	}

	//Non-existent data
	for _, u := range testValidData {
		_, err = prps.Pool.Exec(ctx, "delete from users where login=$1 ", u.Login)

		user, err = prps.GetUserByLogin(ctx, u.Login)
		require.Error(t, err, "get by login error")
	}
}

func TestPUser_UpdateUser(t *testing.T) {
	ctx := context.Background()
	prps := NewPRepository(db)

	var user *model.User
	var err error
	for _, u := range testValidData {
		_, err = prps.Pool.Exec(ctx, "delete from users where login=$1 ", u.Login)
		_, err = prps.CreateUser(ctx, &u)
		require.NoError(t, err, "create error")

		u.Name = "Update"
		err = prps.UpdateUser(ctx, u.Login, &u)
		require.NoError(t, err, "update error")

		user, err = prps.GetUserByLogin(ctx, u.Login)
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
	prps := NewPRepository(db)

	var user *model.User
	var err error
	for _, u := range testValidData {
		_, err = prps.Pool.Exec(ctx, "delete from users where login=$1 ", u.Login)
		_, err = prps.CreateUser(ctx, &u)
		require.NoError(t, err, "create error")

		token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJsb2dpbiI6InRlc3QxIiwiZXhwIjoxNjc0ODMxODE2fQ.jlD1_wrfdK8XjMut236sQDb7B7EOvVjflGZnNUS5o2g"
		err = prps.RefreshUser(ctx, u.Login, token)
		require.NoError(t, err, "refresh error")

		user, err = prps.GetUserByLogin(ctx, u.Login)
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
	prps := NewPRepository(db)

	var err error
	for _, u := range testValidData {
		_, err = prps.Pool.Exec(ctx, "delete from users where login=$1 ", u.Login)
		_, err = prps.CreateUser(ctx, &u)
		require.NoError(t, err, "create error")

		err = prps.DeleteUser(ctx, u.Login)
		require.NoError(t, err, "delete error")

		_, err = prps.GetUserByLogin(ctx, u.Login)
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
