package service

import (
	"GolangInternship/FMicroserviceGRPC/internal/model"
	"context"
)

type Stream interface {
	ProduceUser(ctx context.Context, user *model.User) error
}
