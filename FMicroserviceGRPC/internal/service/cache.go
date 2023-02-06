package service

import (
	"GolangInternship/FMicroserviceGRPC/internal/model"
	"context"
)

type Cache interface {
	GetByLogin(ctx context.Context, login string) (*model.User, bool, error)
	CreateUser(ctx context.Context, user *model.User) error
}
