package service

import (
	"github.com/OVantsevich/GolangInternship/FMicroservice/internal/domain"
	"github.com/OVantsevich/GolangInternship/FMicroservice/internal/repository"
	"github.com/labstack/echo/v4"
	"net/http"
	"unicode"
)

func CreateEntity(c echo.Context) (err error) {
	entity := &domain.Entity{}
	if err = c.Bind(entity); err != nil {
		return
	}

	if !NameValid(entity.Name) || entity.Age < 1 || entity.Age > 100 {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "invalid name or age"}
	}

	if err = repository.Repos.CreateEntity(c.Request().Context(), entity); err != nil {
		switch err.Error() {
		case "database not responding":
			return &echo.HTTPError{Code: http.StatusServiceUnavailable, Message: "server is temporarily unavailable"}
		case "entity with this name already exist":
			return &echo.HTTPError{Code: http.StatusBadRequest, Message: err.Error()}
		}
	}

	return c.JSON(http.StatusCreated, entity)
}

func FindEntity(c echo.Context) (err error) {
	entity := &domain.Entity{}
	if err = c.Bind(entity); err != nil {
		return
	}

	if entity, err = repository.Repos.GetEntityByName(c.Request().Context(), entity.Name); entity == nil {
		switch err.Error() {
		case "database not responding":
			return &echo.HTTPError{Code: http.StatusServiceUnavailable, Message: "server is temporarily unavailable"}
		case "entity with this name doesn't exist":
			return &echo.HTTPError{Code: http.StatusBadRequest, Message: err.Error()}
		}
	}

	return c.JSON(http.StatusOK, entity)
}

func UpdateEntity(c echo.Context) (err error) {
	entity := &domain.Entity{}
	if err = c.Bind(entity); err != nil {
		return
	}

	if !NameValid(entity.Name) || entity.Age < 1 || entity.Age > 100 {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "invalid name or age"}
	}

	if err = repository.Repos.UpdateEntity(c.Request().Context(), entity.Name, entity); err != nil {
		switch err.Error() {
		case "database not responding":
			return &echo.HTTPError{Code: http.StatusServiceUnavailable, Message: "server is temporarily unavailable"}
		case "entity with this name doesn't exist":
			return &echo.HTTPError{Code: http.StatusBadRequest, Message: err.Error()}
		}
	}

	return c.JSON(http.StatusOK, entity)
}

func DeleteEntity(c echo.Context) (err error) {
	entity := &domain.Entity{}
	if err = c.Bind(entity); err != nil {
		return
	}

	if err = repository.Repos.DeleteEntity(c.Request().Context(), entity.Name); err != nil {
		switch err.Error() {
		case "database not responding":
			return &echo.HTTPError{Code: http.StatusServiceUnavailable, Message: "server is temporarily unavailable"}
		case "entity with this name doesn't exist":
			return &echo.HTTPError{Code: http.StatusBadRequest, Message: err.Error()}
		}
	}

	return c.JSON(http.StatusOK, entity)
}

func NameValid(name string) bool {
	if name == "" {
		return false
	}
	for _, r := range name {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}
