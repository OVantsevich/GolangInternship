package repository

import (
	"context"
	"github.com/OVantsevich/GolangInternship/FMicroservice/internal/model"
)

type Stream interface {
	CreatingUser(ctx context.Context, user *model.User) error
}
