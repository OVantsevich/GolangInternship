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

func (r *MRepository) CreateEntity(ctx context.Context, e *Entity) error {
	db := r.Client.Database("entity")
	e.ID = uuid.New().String()
	_, err := db.Collection("entity").InsertOne(ctx, e)
	if err != nil {
		return fmt.Errorf("repository - MRepository - CreateEntity: %v", err)
	}

	return nil
}

func (r *MRepository) GetEntityByName(ctx context.Context, name string) (*Entity, error) {
	e := Entity{}

	db := r.Client.Database("entity")
	result := db.Collection("entity").FindOne(ctx, bson.D{{"name", name}})
	err := result.Decode(&e)
	if err != nil {
		return nil, fmt.Errorf("repository - MRepository - GetEntityByName: %v", err)
	}

	return &e, nil
}
func (r *MRepository) UpdateEntity(ctx context.Context, name string, e *Entity) error {
	db := r.Client.Database("entity")

	filter := bson.D{{"name", name}}
	update := bson.D{{"$set", bson.D{{"age", e.Age}}}}

	_, err := db.Collection("entity").UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("repository - MRepository - UpdateEntity: %v", err)
	}

	return nil
}
func (r *MRepository) DeleteEntity(ctx context.Context, name string) error {
	db := r.Client.Database("entity")

	filter := bson.D{{"name", name}}
	update := bson.D{{"$set", bson.D{{"deleted", true}}}}

	_, err := db.Collection("entity").UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("repository - MRepository - DeleteEntity: %v", err)
	}

	return nil
}
