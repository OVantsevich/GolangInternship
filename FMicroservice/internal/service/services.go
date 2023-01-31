package service

import (
	"context"
	"github.com/OVantsevich/GolangInternship/FMicroservice/internal/model"
)

type User interface {
	Signup(ctx context.Context, user *model.User) (string, string, *model.User, error)
	Login(ctx context.Context, login, password string) (string, string, error)
	Refresh(ctx context.Context, login, userRefreshToken string) (string, string, error)
	Update(ctx context.Context, login string, user *model.User) error
	Delete(ctx context.Context, login string) error
}
