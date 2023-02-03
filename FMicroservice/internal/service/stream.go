package service

import (
	"context"
	"github.com/OVantsevich/GolangInternship/FMicroservice/internal/model"
)

type Stream interface {
	ProduceUser(ctx context.Context, user *model.User) error
}
