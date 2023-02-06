package handler

import (
	"GolangInternship/FMicroservice/internal/model"
	"context"
)

//go:generate mockery --name=UserService --case=underscore --output=./mocks
type UserService interface {
	Signup(ctx context.Context, user *model.User) (string, string, *model.User, error)
	Login(ctx context.Context, login, password string) (string, string, error)
	Refresh(ctx context.Context, login, userRefreshToken string) (string, string, error)
	Update(ctx context.Context, login string, user *model.User) error
	Delete(ctx context.Context, login string) error

	GetByLogin(ctx context.Context, login string) (*model.User, error)
}
