package server

import (
	"github.com/OVantsevich/GolangInternship/FMicroservice/internal/service"
	"github.com/labstack/echo/v4"
)

func SetRoutes(e *echo.Echo) {
	e.POST("/users", service.CreateEntity)
	e.GET("/users", service.FindEntity)
	e.PUT("/users", service.UpdateEntity)
	e.DELETE("/users", service.DeleteEntity)
}
