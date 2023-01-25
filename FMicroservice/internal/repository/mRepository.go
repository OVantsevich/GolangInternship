package repository

import (
	"context"
	"fmt"
	. "github.com/OVantsevich/GolangInternship/FMicroservice/internal/domain"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	MRepository struct {
		Client *mongo.Client
	}
)

func (r *MRepository) CreateUser(ctx context.Context, e *User) error {
	db := r.Client.Database("User")
	e.ID = uuid.New().String()
	_, err := db.Collection("User").InsertOne(ctx, e)
	if err != nil {
		return fmt.Errorf("repository - MRepository - CreateUser: %v", err)
	}

	return nil
}

func (r *MRepository) GetUserByName(ctx context.Context, name string) (*User, error) {
	e := User{}

	db := r.Client.Database("User")
	result := db.Collection("User").FindOne(ctx, bson.D{{"name", name}})
	err := result.Decode(&e)
	if err != nil {
		return nil, fmt.Errorf("repository - MRepository - GetUserByName: %v", err)
	}

	return &e, nil
}
func (r *MRepository) UpdateUser(ctx context.Context, name string, e *User) error {
	db := r.Client.Database("User")

	filter := bson.D{{"name", name}}
	update := bson.D{{"$set", bson.D{{"age", e.Age}}}}

	_, err := db.Collection("User").UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("repository - MRepository - UpdateUser: %v", err)
	}

	return nil
}
func (r *MRepository) DeleteUser(ctx context.Context, name string) error {
	db := r.Client.Database("User")

	filter := bson.D{{"name", name}}
	update := bson.D{{"$set", bson.D{{"deleted", true}}}}

	_, err := db.Collection("User").UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("repository - MRepository - DeleteUser: %v", err)
	}

	return nil
}
