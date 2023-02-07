// Package repository mUser
package repository

import (
	"GolangInternship/FMicroserviceGRPC/internal/model"
	"context"
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

// MUser mongo entity
type MUser struct {
	Client *mongo.Client
}

type mongoUser struct {
	*model.User `bson:"user"`
	Role        string `bson:"Role"`
	Deleted     bool   `bson:"Deleted"`
}

// CreateUser create user
func (r *MUser) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	collection := r.Client.Database("userService").Collection("users")
	ID := uuid.New().String()
	_, err := collection.InsertOne(ctx, mongoUser{
		User: &model.User{
			ID:       ID,
			Login:    user.Login,
			Email:    user.Email,
			Password: user.Password,
			Name:     user.Name,
			Age:      user.Age,
			Token:    "",
			Created:  time.Now(),
			Updated:  time.Now(),
		},
		Role:    "user",
		Deleted: false})
	if err != nil {
		return nil, fmt.Errorf("MUser - CreateUser - InsertOne: %w", err)
	}

	return user, nil
}

// GetUserByLogin get user by login
func (r *MUser) GetUserByLogin(ctx context.Context, login string) (*model.User, error) {
	user := mongoUser{}

	collection := r.Client.Database("userService").Collection("users")
	result := collection.FindOne(ctx, bson.D{primitive.E{Key: "user.login", Value: login}, primitive.E{Key: "Deleted", Value: false}})
	err := result.Decode(&user)
	if err != nil {
		return nil, fmt.Errorf("MUser - GetUserByName - Decode: %w", err)
	}

	return user.User, nil
}

// UpdateUser update user
func (r *MUser) UpdateUser(ctx context.Context, login string, user *model.User) error {
	collection := r.Client.Database("userService").Collection("users")

	filter := bson.D{primitive.E{Key: "user.login", Value: login}, primitive.E{Key: "Deleted", Value: false}}
	update := bson.D{primitive.E{Key: "$set", Value: bson.D{
		primitive.E{Key: "user.email", Value: user.Email},
		primitive.E{Key: "user.name", Value: user.Name},
		primitive.E{Key: "user.age", Value: user.Age},
		primitive.E{Key: "user.update", Value: user.Updated}}}}

	userResult := model.User{}
	err := collection.FindOneAndUpdate(ctx, filter, update).Decode(&userResult)
	if err != nil {
		return fmt.Errorf("MUser - UpdateUser - UpdateOne: %w", err)
	}

	return nil
}

// RefreshUser refresh user
func (r *MUser) RefreshUser(ctx context.Context, login, token string) error {
	collection := r.Client.Database("userService").Collection("users")

	filter := bson.D{primitive.E{Key: "user.login", Value: login}, primitive.E{Key: "Deleted", Value: false}}
	update := bson.D{primitive.E{Key: "$set", Value: bson.D{
		primitive.E{Key: "user.token", Value: token},
		primitive.E{Key: "user.updated", Value: time.Now()}}}}

	userResult := model.User{}
	err := collection.FindOneAndUpdate(ctx, filter, update).Decode(&userResult)
	if err != nil {
		return fmt.Errorf("MUser - RefreshUser - UpdateOne: %w", err)
	}

	return nil
}

// DeleteUser delete user
func (r *MUser) DeleteUser(ctx context.Context, login string) error {
	collection := r.Client.Database("userService").Collection("users")

	filter := bson.D{primitive.E{Key: "user.login", Value: login}, primitive.E{Key: "Deleted", Value: false}}
	update := bson.D{primitive.E{Key: "$set", Value: bson.D{
		primitive.E{Key: "Deleted", Value: true},
		primitive.E{Key: "user.updated", Value: time.Now()}}}}

	userResult := model.User{}
	err := collection.FindOneAndUpdate(ctx, filter, update).Decode(&userResult)
	if err != nil {
		return fmt.Errorf("MUser - DeleteUser - UpdateOne: %w", err)
	}

	return nil
}
