package service

import (
	"context"
	. "github.com/OVantsevich/GolangInternship/FMicroservice/internal/domain"
	. "github.com/OVantsevich/GolangInternship/FMicroservice/internal/repository"
	"github.com/labstack/echo/v4"
	passwordvalidator "github.com/wagslane/go-password-validator"
	"net/http"
)

type UserService struct {
	rps Repository
}

func NewUserService(rps *Repository) *UserService {
	return &UserService{rps: *rps}
}

func (es *UserService) CreateUser(ctx context.Context, e *User) (err error) {

	if err = passwordvalidator.Validate(e.Password, 50); err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: err.Error()}
	}

	if err = es.rps.CreateUser(ctx, e); err != nil {
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	return
}

func (es *UserService) FindUser(ctx context.Context, name string) (e *User, err error) {
	if e, err = es.rps.GetUserByName(ctx, name); err != nil {
		return nil, &echo.HTTPError{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	return
}

func (es *UserService) UpdateUser(ctx context.Context, e *User) (err error) {
	if err = es.rps.UpdateUser(ctx, e.Name, e); err != nil {
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	return
}

func (es *UserService) DeleteUser(ctx context.Context, name string) (err error) {
	if err = es.rps.DeleteUser(ctx, name); err != nil {
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	return
}

//func ErrorDB(err error) *echo.HTTPError {
//	switch err.Error() {
//	case "database not responding":
//		return &echo.HTTPError{Code: http.StatusServiceUnavailable, Message: "handler is temporarily unavailable"}
//	case "User with this name already exist":
//		return &echo.HTTPError{Code: http.StatusBadRequest, Message: err.Error()}
//	default:
//		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: err.Error()}
//	}
//}
