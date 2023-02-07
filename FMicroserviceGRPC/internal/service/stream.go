package service

import (
	"GolangInternship/FMicroserviceGRPC/internal/model"
	"context"
)

//go:generate mockery --name=Stream --case=underscore --output=./mocks
type Stream interface {
	ProduceUser(ctx context.Context, user *model.User) error
}
