package repository

import (
	"context"
	"github.com/OVantsevich/GolangInternship/FMicroservice/internal/model"
)

type Cache interface {
	GetByLogin(ctx context.Context, login string) (*model.User, error)
	CreateUser(ctx context.Context, user *model.User) error
}
