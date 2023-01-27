package repository

import (
	"context"
	"fmt"
	. "github.com/OVantsevich/GolangInternship/FMicroservice/internal/domain"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type (
	MRepository struct {
		Client *mongo.Client
	}
)

func (r *MRepository) CreateUser(ctx context.Context, user *User) error {
	collection := r.Client.Database("userService").Collection("users")
	ID := uuid.New().String()
	_, err := collection.InsertOne(ctx, User{
		ID:       ID,
		Login:    user.Login,
		Email:    user.Email,
		Password: user.Password,
		Name:     user.Name,
		Age:      user.Age,
		Token:    "",
		Deleted:  false,
		Created:  time.Now(),
		Updated:  time.Now(),
	})
	if err != nil {
		return fmt.Errorf("repository - MRepository - CreateUser: %v", err)
	}

	return nil
}

func (r *MRepository) GetUserByLogin(ctx context.Context, login string) (*User, error) {
	user := User{}

	collection := r.Client.Database("userService").Collection("users")
	result := collection.FindOne(ctx, bson.D{{"login", login}})
	err := result.Decode(&user)
	if err != nil {
		return nil, fmt.Errorf("repository - MRepository - GetUserByName: %v", err)
	}

	return &user, nil
}
func (r *MRepository) UpdateUser(ctx context.Context, login string, user *User) error {
	collection := r.Client.Database("userService").Collection("users")

	filter := bson.D{{"login", login}}
	update := bson.D{{"$set", bson.D{
		{"email", user.Email},
		{"name", user.Name},
		{"age", user.Age},
		{"update", user.Updated}}}}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("repository - MRepository - UpdateUser: %v", err)
	}

	return nil
}

func (r *MRepository) RefreshUser(ctx context.Context, login, token string) error {
	collection := r.Client.Database("userService").Collection("users")

	filter := bson.D{{"login", login}}
	update := bson.D{{"$set", bson.D{
		{"token", token},
		{"updated", time.Now()}}}}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("repository - MRepository - UpdateUser: %v", err)
	}

	return nil
}

func (r *MRepository) DeleteUser(ctx context.Context, login string) error {
	collection := r.Client.Database("userService").Collection("users")

	filter := bson.D{{"login", login}}
	update := bson.D{{"$set", bson.D{
		{"deleted", true},
		{"updated", time.Now()}}}}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("repository - MRepository - DeleteUser: %v", err)
	}

	return nil
}
