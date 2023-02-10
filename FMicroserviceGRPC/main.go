// Package main Main package
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
	"github.com/golang-jwt/jwt/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"net"
)

func jwtAuth(keyFunc func(token *jwt.Token) (interface{}, error)) grpc.ServerOption {
	return grpc.UnaryInterceptor(func(ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {

		if info.FullMethod != "/UserService/Signup" && info.FullMethod != "/UserService/Login" && info.FullMethod != "/UserService/Refresh" {
			md, ok := metadata.FromIncomingContext(ctx)
			if !ok {
				return nil, status.Errorf(codes.InvalidArgument, "Retrieving metadata is failed")
			}

			authHeader, ok := md["authorization"]
			if !ok {
				return nil, status.Errorf(codes.Unauthenticated, "Authorization token is not supplied")
			}

			token := authHeader[0]

			claims, err := verify(token, keyFunc)
			if err != nil {
				return nil, err
			}
			ctx = context.WithValue(ctx, "user", claims)
		}

		h, err := handler(ctx, req)

		return h, err
	})
}

func verify(token string, keyFunc func(token *jwt.Token) (interface{}, error)) (claims *service.CustomClaims, err error) {
	claims = &service.CustomClaims{}

	_, err = jwt.ParseWithClaims(
		token,
		claims,
		keyFunc,
	)
	if err != nil {
		err = fmt.Errorf("invalid token: %w", err)
	}

	return
}

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

	rds.ConsumeUser("example")

	userService := service.NewUserServiceClassic(repos, rds, rds, cfg.JwtKey)

	ns := grpc.NewServer(jwtAuth(func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JwtKey), nil
	}))
	fileService := service.NewFile("fileStore")
	server := handler.NewUserHandlerClassic(userService, fileService, cfg.JwtKey)
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
