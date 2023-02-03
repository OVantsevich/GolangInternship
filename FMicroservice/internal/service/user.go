package service

import (
	"context"
	"github.com/OVantsevich/GolangInternship/FMicroservice/internal/model"
)

type User interface {
	CreateUser(ctx context.Context, user *model.User) (*model.User, error)
	GetUserByLogin(ctx context.Context, login string) (*model.User, error)
	UpdateUser(ctx context.Context, login string, user *model.User) error
	RefreshUser(ctx context.Context, login, token string) error
	DeleteUser(ctx context.Context, login string) error
}
