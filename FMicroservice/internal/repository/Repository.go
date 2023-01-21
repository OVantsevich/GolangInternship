package repository

import (
	"context"
	"github.com/OVantsevich/GolangInternship/FMicroservice/internal/domain"
)

type Repository interface {
	OpenPool(ctx context.Context) error
	ClosePool()
	CreateEntity(ctx context.Context, e *domain.Entity) error
	GetEntityByName(ctx context.Context, name string) (*domain.Entity, error)
	UpdateEntity(ctx context.Context, name string, e *domain.Entity) error
	DeleteEntity(ctx context.Context, name string) error
}

var Repos Repository
