package main

import (
	"context"
	"fmt"
	"github.com/OVantsevich/GolangInternship/FMicroservice/internal/config"
	"github.com/OVantsevich/GolangInternship/FMicroservice/internal/handler"
	"github.com/OVantsevich/GolangInternship/FMicroservice/internal/repository"
	"github.com/OVantsevich/GolangInternship/FMicroservice/internal/service"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io"
	"net/http"
	"os"

	_ "github.com/OVantsevich/GolangInternship/FMicroservice/docs"
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

// @title Swagger Example API
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:12345
// @BasePath /

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	e := echo.New()

	cfg, err := config.NewConfig()
	if err != nil {
		logrus.Fatal(err)
	}

	e.Use(echojwt.WithConfig(echojwt.Config{
		Skipper: func(c echo.Context) bool {
			if c.Path() == "/login" || c.Path() == "/signup" || c.Path() == "/swagger/*" {
				return true
			}
			return false
		},
		KeyFunc: func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.JwtKey), nil
		},
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(service.CustomClaims)
		},
	}))

	var repos repository.User
	repos, err = DBConnection(cfg)
	if err != nil {
		logrus.Fatal(err)
	}
	defer ClosePool(cfg, repos)

	userService := service.NewUserServiceClassic(repos, cfg.JwtKey)
	userHandler := handler.NewUserHandlerClassic(userService)

	e.Validator = &CustomValidator{validator: validator.New()}

	e.POST("/signup", userHandler.Signup)
	e.POST("/login", userHandler.Login)
	e.PUT("/update", userHandler.Update)
	e.DELETE("/delete", userHandler.Delete)
	e.GET("/refresh", userHandler.Refresh)

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.GET("/", func(c echo.Context) error {
		return c.File("index.html")
	})
	e.GET("/file", func(c echo.Context) error {
		return c.File("file.svg")
	})
	e.POST("/upload", upload)

	e.Logger.Fatal(e.Start(":12345"))
}

func DBConnection(Cfg *config.Config) (repository.User, error) {
	switch Cfg.CurrentDB {
	case "postgres":
		pool, err := pgxpool.New(context.Background(), Cfg.PostgresUrl)
		if err != nil {
			return nil, fmt.Errorf("invalid configuration data: %v", err)
		}
		if err = pool.Ping(context.Background()); err != nil {
			return nil, fmt.Errorf("database not responding: %v", err)
		}
		return &repository.PUser{Pool: pool}, nil
	case "mongo":
		client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(Cfg.MongoURL))
		if err != nil {
			return nil, fmt.Errorf("mongoDB connection: %v", err)
		}
		err = client.Ping(context.Background(), nil)
		if err != nil {
			return nil, fmt.Errorf("database not responding: %v", err)
		}
		return &repository.MUser{Client: client}, nil
	}
	return nil, nil
}

func ClosePool(Cfg *config.Config, r interface{}) {
	switch Cfg.CurrentDB {
	case "postgres":
		pr := r.(repository.PUser)
		if pr.Pool != nil {
			pr.Pool.Close()
		}
	case "mongo":
		pr := r.(repository.MUser)
		if pr.Client != nil {
			err := pr.Client.Disconnect(context.Background())
			if err != nil {
				logrus.Fatalf("mongoDB disconnecting: %v", err)
			}
		}
	}
}
