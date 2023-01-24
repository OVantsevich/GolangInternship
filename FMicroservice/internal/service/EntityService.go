package service

import (
	"context"
	. "github.com/OVantsevich/GolangInternship/FMicroservice/internal/domain"
	. "github.com/OVantsevich/GolangInternship/FMicroservice/internal/repository"
	"github.com/labstack/echo/v4"
	"net/http"
)

type EntityService struct {
	rps Repository
}

func NewEntityService(rps *Repository) *EntityService {
	return &EntityService{rps: *rps}
}

func (es *EntityService) CreateEntity(ctx context.Context, e *Entity) (err error) {
	if err = es.rps.CreateEntity(ctx, e); err != nil {
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	return
}

func (es *EntityService) FindEntity(ctx context.Context, name string) (e *Entity, err error) {
	if e, err = es.rps.GetEntityByName(ctx, name); err != nil {
		return nil, &echo.HTTPError{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	return
}

func (es *EntityService) UpdateEntity(ctx context.Context, e *Entity) (err error) {
	if err = es.rps.UpdateEntity(ctx, e.Name, e); err != nil {
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	return
}

func (es *EntityService) DeleteEntity(ctx context.Context, name string) (err error) {
	if err = es.rps.DeleteEntity(ctx, name); err != nil {
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	return
}

//func ErrorDB(err error) *echo.HTTPError {
//	switch err.Error() {
//	case "database not responding":
//		return &echo.HTTPError{Code: http.StatusServiceUnavailable, Message: "handler is temporarily unavailable"}
//	case "entity with this name already exist":
//		return &echo.HTTPError{Code: http.StatusBadRequest, Message: err.Error()}
//	default:
//		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: err.Error()}
//	}
//}
