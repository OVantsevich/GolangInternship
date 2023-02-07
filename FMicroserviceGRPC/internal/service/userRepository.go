package service

import (
	"GolangInternship/FMicroserviceGRPC/internal/model"
	"context"
)

//go:generate mockery --name=UserRepository --case=underscore --output=./mocks
type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) (*model.User, error)
	GetUserByLogin(ctx context.Context, login string) (*model.User, error)
	UpdateUser(ctx context.Context, login string, user *model.User) error
	RefreshUser(ctx context.Context, login, token string) error
	DeleteUser(ctx context.Context, login string) error
}
