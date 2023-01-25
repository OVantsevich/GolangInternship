package repository

import (
	"context"
	. "github.com/OVantsevich/GolangInternship/FMicroservice/internal/domain"
)

type Repository interface {
	CreateUser(ctx context.Context, e *User) error
	GetUserByName(ctx context.Context, name string) (*User, error)
	UpdateUser(ctx context.Context, name string, e *User) error
	DeleteUser(ctx context.Context, name string) error
}
