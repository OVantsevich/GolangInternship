// Package main Main package
package main

import (
	_ "GolangInternship/FMicroservice/docs"
	"GolangInternship/FMicroservice/internal/config"
	"GolangInternship/FMicroservice/internal/handler"
	"GolangInternship/FMicroservice/internal/repository"
	"GolangInternship/FMicroservice/internal/service"
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
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

	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(file.Filename)
	if err != nil {
		return err
	}
	defer dst.Close()

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

	var repos service.UserClassicRepository
	repos, err = DBConnection(cfg)
	if err != nil {
		logrus.Fatal(err)
	}
	defer ClosePool(cfg, repos)

	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer client.Close()

	rds := &repository.Redis{Client: *client}

	rds.RedisStreamInit(context.Background())
	rds.ConsumeUser("example")

	userService := service.NewUserServiceClassic(repos, rds, rds, cfg.JwtKey)
	userHandler := handler.NewUserHandlerClassic(userService)

	e.Validator = &CustomValidator{validator: validator.New()}

	e.POST("/signup", userHandler.Signup)
	e.POST("/login", userHandler.Login)
	e.PUT("/update", userHandler.Update)
	e.DELETE("/delete", userHandler.Delete)
	e.GET("/refresh", userHandler.Refresh)

	admin := e.Group("/admin")
	admin.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user := c.Get("user").(*jwt.Token)
			claims := user.Claims
			if claims.(*service.CustomClaims).Role == "admin" {
				return next(c)
			}
			return echo.NewHTTPError(http.StatusForbidden, "access denied")
		}
	})
	admin.GET("/userByLogin", userHandler.UserByLogin)

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

func DBConnection(Cfg *config.Config) (service.UserClassicRepository, error) {
	switch Cfg.CurrentDB {
	case "postgres":
		pool, err := pgxpool.New(context.Background(), Cfg.PostgresURL)
		if err != nil {
			return nil, fmt.Errorf("invalid configuration data: %v", err)
		}
		if err = pool.Ping(context.Background()); err != nil {
			return nil, fmt.Errorf("database not responding: %v", err)
		}
		return repository.NewPostgresRepository(pool), nil
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
