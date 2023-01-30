package repository

import (
	"context"
	"fmt"
	"github.com/OVantsevich/GolangInternship/FMicroservice/internal/model"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type MUser struct {
	Client *mongo.Client
}

func (r *MUser) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	collection := r.Client.Database("userService").Collection("users")
	ID := uuid.New().String()
	_, err := collection.InsertOne(ctx, model.User{
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
		return nil, fmt.Errorf("MUser - CreateUser - InsertOne: %w", err)
	}

	return user, nil
}

func (r *MUser) GetUserByLogin(ctx context.Context, login string) (*model.User, error) {
	user := model.User{}

	collection := r.Client.Database("userService").Collection("users")
	result := collection.FindOne(ctx, bson.D{{"login", login}, {"deleted", false}})
	err := result.Decode(&user)
	if err != nil {
		return nil, fmt.Errorf("MUser - GetUserByName - Decode: %w", err)
	}

	return &user, nil
}
func (r *MUser) UpdateUser(ctx context.Context, login string, user *model.User) error {
	collection := r.Client.Database("userService").Collection("users")

	filter := bson.D{{"login", login}, {"deleted", false}}
	update := bson.D{{"$set", bson.D{
		{"email", user.Email},
		{"name", user.Name},
		{"age", user.Age},
		{"update", user.Updated}}}}

	userResult := model.User{}
	err := collection.FindOneAndUpdate(ctx, filter, update).Decode(&userResult)
	if err != nil {
		return fmt.Errorf("MUser - UpdateUser - UpdateOne: %w", err)
	}

	return nil
}

func (r *MUser) RefreshUser(ctx context.Context, login, token string) error {
	collection := r.Client.Database("userService").Collection("users")

	filter := bson.D{{"login", login}, {"deleted", false}}
	update := bson.D{{"$set", bson.D{
		{"token", token},
		{"updated", time.Now()}}}}

	userResult := model.User{}
	err := collection.FindOneAndUpdate(ctx, filter, update).Decode(&userResult)
	if err != nil {
		return fmt.Errorf("MUser - RefreshUser - UpdateOne: %w", err)
	}

	return nil
}

func (r *MUser) DeleteUser(ctx context.Context, login string) error {
	collection := r.Client.Database("userService").Collection("users")

	filter := bson.D{{"login", login}, {"deleted", false}}
	update := bson.D{{"$set", bson.D{
		{"deleted", true},
		{"updated", time.Now()}}}}

	userResult := model.User{}
	err := collection.FindOneAndUpdate(ctx, filter, update).Decode(&userResult)
	if err != nil {
		return fmt.Errorf("MUser - DeleteUser - UpdateOne: %w", err)
	}

	return nil
}
