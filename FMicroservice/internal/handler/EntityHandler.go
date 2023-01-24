package handler

import (
	. "github.com/OVantsevich/GolangInternship/FMicroservice/internal/domain"
	. "github.com/OVantsevich/GolangInternship/FMicroservice/internal/service"
	"github.com/labstack/echo/v4"
	"net/http"
)

type EntityHandler struct {
	es EntityService
}

func NewEntityHandler(es *EntityService) *EntityHandler {
	return &EntityHandler{es: *es}
}

func (eh *EntityHandler) CreateEntity(c echo.Context) (err error) {
	entity := &Entity{}
	if err = c.Bind(entity); err != nil {
		return
	}

	if err = c.Validate(entity); err != nil {
		return
	}

	if err = eh.es.CreateEntity(c.Request().Context(), entity); err != nil {
		return
	}

	return c.JSON(http.StatusCreated, entity)
}

func (eh *EntityHandler) GetEntityByName(c echo.Context) (err error) {
	entity := &Entity{}
	if err = c.Bind(entity); err != nil {
		return
	}

	if entity, err = eh.es.FindEntity(c.Request().Context(), entity.Name); err != nil {
		return
	}

	return c.JSON(http.StatusOK, entity)
}

func (eh *EntityHandler) UpdateEntity(c echo.Context) (err error) {
	entity := &Entity{}
	if err = c.Bind(entity); err != nil {
		return
	}

	if err = c.Validate(entity); err != nil {
		return
	}

	if err = eh.es.UpdateEntity(c.Request().Context(), entity); err != nil {
		return
	}

	return c.JSON(http.StatusOK, entity)
}

func (eh *EntityHandler) DeleteEntity(c echo.Context) (err error) {
	entity := &Entity{}
	if err = c.Bind(entity); err != nil {
		return
	}

	if err = eh.es.DeleteEntity(c.Request().Context(), entity.Name); err != nil {
		return
	}

	return c.JSON(http.StatusOK, entity)
}
