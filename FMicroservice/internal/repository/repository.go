package repository

import (
	"context"
	. "github.com/OVantsevich/GolangInternship/FMicroservice/internal/domain"
)

type Repository interface {
	CreateUser(ctx context.Context, user *User) error
	GetUserByLogin(ctx context.Context, login string) (*User, error)
	UpdateUser(ctx context.Context, login string, user *User) error
	RefreshUser(ctx context.Context, login, token string) error
	DeleteUser(ctx context.Context, login string) error
}
