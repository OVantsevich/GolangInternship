package repository

import (
	"context"
	. "github.com/OVantsevich/GolangInternship/FMicroservice/internal/domain"
)

type Repository interface {
	CreateEntity(ctx context.Context, e *Entity) error
	GetEntityByName(ctx context.Context, name string) (*Entity, error)
	UpdateEntity(ctx context.Context, name string, e *Entity) error
	DeleteEntity(ctx context.Context, name string) error
}
