package main

import (
	"GolangInternship/FMicroserviceGRPC/internal/config"
	"GolangInternship/FMicroserviceGRPC/internal/handler"
	"GolangInternship/FMicroserviceGRPC/internal/repository"
	"GolangInternship/FMicroserviceGRPC/internal/service"
	pr "GolangInternship/FMicroserviceGRPC/proto"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"net"
)

func main() {
	listen, err := net.Listen("tcp", "localhost:12344")
	if err != nil {
		defer logrus.Fatalf("error while listening port: %e", err)
	}
	fmt.Println("Server successfully started on port :12344...")
	cfg, err := config.NewConfig()
	if err != nil {
		logrus.Fatal(err)
	}

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

	ns := grpc.NewServer()
	server := handler.NewUserHandlerClassic(userService, cfg.JwtKey)
	pr.RegisterUserServiceServer(ns, server)

	if err = ns.Serve(listen); err != nil {
		defer logrus.Fatalf("error while listening server: %e", err)
	}
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
