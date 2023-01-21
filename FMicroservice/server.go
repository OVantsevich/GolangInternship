package main

import (
	"context"
	. "github.com/OVantsevich/GolangInternship/FMicroservice/internal/domain"
	. "github.com/OVantsevich/GolangInternship/FMicroservice/internal/repository"
	. "github.com/OVantsevich/GolangInternship/FMicroservice/internal/server"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
	"os"
)

func main() {
	InitConfig()
	Repos = &PRepository{}
	Repos.OpenPool(context.Background())
	defer Repos.ClosePool()

	e := echo.New()

	var logger = log.New()
	logger.Out = os.Stdout
	log.SetReportCaller(true)
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, values middleware.RequestLoggerValues) error {
			log.WithFields(log.Fields{
				"URI":    values.URI,
				"status": values.Status,
			}).Info("request")

			return nil
		},
	}))

	//e.Use(echojwt.WithConfig(echojwt.Config{
	//	SigningKey: []byte(domain.Cfg.JwtKey),
	//}))

	SetRoutes(e)

	e.Logger.Fatal(e.Start(":12345"))
}
