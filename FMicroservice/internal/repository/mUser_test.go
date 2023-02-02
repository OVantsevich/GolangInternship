package repository

import (
	"context"
	"github.com/OVantsevich/GolangInternship/FMicroservice/internal/model"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
)

var mrps *MUser

var mTestValidData = []MongoUser{
	{
		User: &model.User{
			Name:     `NAME`,
			Age:      1,
			Login:    `CreateLOGIN11`,
			Email:    `LOGIN1@gmail.com`,
			Password: `LOGIN123456789`,
		},
		role: "user",
	},
	{
		User: &model.User{
			Name:     `NAME`,
			Age:      1,
			Login:    `CreateLOGIN22`,
			Email:    `LOGIN2@gmail.com`,
			Password: `PASSWORD123456789`,
		},
		role: "user",
	},
}
var mTestNoValidData = []MongoUser{
	{
		User: &model.User{
			Name:     `nameEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEE`,
			Age:      22222,
			Login:    `LOGIN2`,
			Email:    `LOGIN2@gmail.com`,
			Password: `PASSWORD123`,
		},
		role: "user",
	},
	{
		User: &model.User{
			Name:     `NAME`,
			Age:      2,
			Login:    `LOGIN1AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA`,
			Email:    `LOGIN1@gmail.com`,
			Password: `LOGIN23102002`,
		},
		role: "user",
	},
}

func NewMRepository(client *mongo.Client) *MUser {
	return &MUser{Client: client}
}

func TestMUser_CreateUser(t *testing.T) {
	ctx := context.Background()
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://mongo:mongo@localhost:27017"))
	require.NoError(t, err, "new client error")
	mrps = NewMRepository(client)

	for _, u := range mTestValidData {
		_, err = mrps.CreateUser(ctx, u.User)
		require.NoError(t, err, "create error")

		_, err = mrps.Client.Database("userService").Collection("users").DeleteOne(
			ctx, bson.D{{"user.login", u.Login}})
	}

	// Already existing data
	for _, u := range mTestValidData {
		_, err = mrps.CreateUser(ctx, u.User)
		require.NoError(t, err, "create error")

		_, err = mrps.CreateUser(ctx, u.User)
		require.Error(t, err, "create error")

		_, err = mrps.Client.Database("userService").Collection("users").DeleteOne(
			ctx, bson.D{{"user.login", u.Login}})
	}
}

func TestMUser_GetUserByLogin(t *testing.T) {
	ctx := context.Background()
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://mongo:mongo@localhost:27017"))
	require.NoError(t, err, "new client error")
	mrps = NewMRepository(client)

	var user *model.User
	for _, u := range mTestValidData {
		_, err = mrps.Client.Database("userService").Collection("users").DeleteOne(
			ctx, bson.D{{"user.login", u.Login}})
		_, err = mrps.CreateUser(ctx, u.User)
		require.NoError(t, err, "create error")

		user, err = mrps.GetUserByLogin(ctx, u.Login)
		require.Equal(t, u.Password, user.Password)
		require.Equal(t, u.Email, user.Email)
		require.NoError(t, err, "get by login error")

		_, err = mrps.Client.Database("userService").Collection("users").DeleteOne(
			ctx, bson.D{{"user.login", u.Login}})
	}

	//Non-existent data
	for _, u := range mTestValidData {
		_, err = mrps.Client.Database("userService").Collection("users").DeleteOne(
			ctx, bson.D{{"user.login", u.Login}})

		user, err = mrps.GetUserByLogin(ctx, u.Login)
		require.Error(t, err, "get by login error")
	}
}

func TestMUser_UpdateUser(t *testing.T) {
	ctx := context.Background()
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://mongo:mongo@localhost:27017"))
	require.NoError(t, err, "new client error")
	mrps = NewMRepository(client)

	var user *model.User
	for _, u := range mTestValidData {
		_, err = mrps.Client.Database("userService").Collection("users").DeleteOne(
			ctx, bson.D{{"user.login", u.Login}})
		_, err = mrps.CreateUser(ctx, u.User)
		require.NoError(t, err, "create error")

		u.Name = "Update"
		err = mrps.UpdateUser(ctx, u.Login, u.User)
		require.NoError(t, err, "update error")

		user, err = mrps.GetUserByLogin(ctx, u.Login)
		require.Equal(t, "Update", user.Name)
		require.NoError(t, err, "get by login error")

		_, err = mrps.Client.Database("userService").Collection("users").DeleteOne(
			ctx, bson.D{{"user.login", u.Login}})
	}

	//Non-existent data
	for _, u := range mTestValidData {
		_, err = mrps.Client.Database("userService").Collection("users").DeleteOne(
			ctx, bson.D{{"user.login", u.Login}})

		err = mrps.UpdateUser(ctx, u.Login, u.User)
		require.Error(t, err, "update error")
	}
}

func TestMUser_RefreshUser(t *testing.T) {
	ctx := context.Background()
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://mongo:mongo@localhost:27017"))
	require.NoError(t, err, "new client error")
	mrps = NewMRepository(client)

	var user *model.User
	for _, u := range mTestValidData {
		_, err = mrps.Client.Database("userService").Collection("users").DeleteOne(
			ctx, bson.D{{"user.login", u.Login}})
		_, err = mrps.CreateUser(ctx, u.User)
		require.NoError(t, err, "create error")

		token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJsb2dpbiI6InRlc3QxIiwiZXhwIjoxNjc0ODMxODE2fQ.jlD1_wrfdK8XjMut236sQDb7B7EOvVjflGZnNUS5o2g"
		err = mrps.RefreshUser(ctx, u.Login, token)
		require.NoError(t, err, "refresh error")

		user, err = mrps.GetUserByLogin(ctx, u.Login)
		require.Equal(t, token, user.Token)
		require.NoError(t, err, "get by login error")

		_, err = mrps.Client.Database("userService").Collection("users").DeleteOne(
			ctx, bson.D{{"user.login", u.Login}})
	}

	//Non-existent data
	for _, u := range mTestValidData {
		_, err = mrps.Client.Database("userService").Collection("users").DeleteOne(
			ctx, bson.D{{"user.login", u.Login}})

		token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJsb2dpbiI6InRlc3QxIiwiZXhwIjoxNjc0ODMxODE2fQ.jlD1_wrfdK8XjMut236sQDb7B7EOvVjflGZnNUS5o2g"
		err = mrps.RefreshUser(ctx, u.Login, token)
		require.Error(t, err, "refresh error")
	}
}

func TestMUser_DeleteUser(t *testing.T) {
	ctx := context.Background()
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://mongo:mongo@localhost:27017"))
	require.NoError(t, err, "new client error")
	mrps = NewMRepository(client)

	for _, u := range mTestValidData {
		_, err = mrps.Client.Database("userService").Collection("users").DeleteOne(
			ctx, bson.D{{"user.login", u.Login}})
		_, err = mrps.CreateUser(ctx, u.User)
		require.NoError(t, err, "create error")

		err = mrps.DeleteUser(ctx, u.Login)
		require.NoError(t, err, "delete error")

		_, err = mrps.GetUserByLogin(ctx, u.Login)
		require.Error(t, err, "get by login error")

		_, err = mrps.Client.Database("userService").Collection("users").DeleteOne(
			ctx, bson.D{{"user.login", u.Login}})
	}

	//Non-existent data
	for _, u := range mTestValidData {
		_, err = mrps.Client.Database("userService").Collection("users").DeleteOne(
			ctx, bson.D{{"user.login", u.Login}})

		err = mrps.DeleteUser(ctx, u.Login)
		require.Error(t, err, "delete error")
	}
}
