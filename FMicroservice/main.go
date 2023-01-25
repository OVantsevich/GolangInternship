package main

import (
	"context"
	"fmt"
	. "github.com/OVantsevich/GolangInternship/FMicroservice/internal/config"
	. "github.com/OVantsevich/GolangInternship/FMicroservice/internal/handler"
	. "github.com/OVantsevich/GolangInternship/FMicroservice/internal/repository"
	. "github.com/OVantsevich/GolangInternship/FMicroservice/internal/service"
	"github.com/go-playground/validator"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io"
	"net/http"
	"os"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func upload(c echo.Context) error {
	name := c.FormValue("name")
	email := c.FormValue("email")

	//-----------
	// Read file
	//-----------

	// Source
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Destination
	dst, err := os.Create(file.Filename)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	return c.HTML(http.StatusOK, fmt.Sprintf("<p>File %s uploaded successfully with fields name=%s and email=%s.</p>", file.Filename, name, email))
}

func main() {
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

	cfg, err := NewConfig()
	if err != nil {
		e.Logger.Fatal(err)
	}

	var repos Repository
	repos, err = DBConnection(cfg)
	if err != nil {
		e.Logger.Fatal(err)
	}
	defer ClosePool(cfg, repos)

	service := NewEntityService(&repos)
	handler := NewEntityHandler(service)

	e.Validator = &CustomValidator{validator: validator.New()}

	e.POST("/entity", handler.CreateEntity)
	e.GET("/entity", handler.GetEntityByName)
	e.PUT("/entity", handler.UpdateEntity)
	e.DELETE("/entity", handler.DeleteEntity)

	e.GET("/", func(c echo.Context) error {
		return c.File("index.html")
	})
	e.GET("/file", func(c echo.Context) error {
		return c.File("file.svg")
	})
	e.POST("/upload", upload)

	e.Logger.Fatal(e.Start(":12345"))
}

func DBConnection(Cfg *Config) (Repository, error) {
	switch Cfg.CurrentDB {
	case "postgres":
		pool, err := pgxpool.New(context.Background(), Cfg.PostgresUrl)
		if err != nil {
			return nil, fmt.Errorf("invalid configuration data: %v", err)
		}
		if err = pool.Ping(context.Background()); err != nil {
			return nil, fmt.Errorf("database not responding: %v", err)
		}
		return &PRepository{Pool: pool}, nil
	case "mongo":
		client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(Cfg.MongoURL))
		if err != nil {
			return nil, fmt.Errorf("mongoDB connection: %v", err)
		}
		return &MRepository{Client: client}, nil
	}
	return nil, nil
}

func ClosePool(Cfg *Config, r interface{}) {
	switch Cfg.CurrentDB {
	case "postgres":
		pr := r.(PRepository)
		if pr.Pool != nil {
			pr.Pool.Close()
		}
	case "mongo":
		pr := r.(MRepository)
		if pr.Client != nil {
			err := pr.Client.Disconnect(context.Background())
			if err != nil {
				log.Fatalf("mongoDB disconnecting: %v", err)
			}
		}
	}
}
