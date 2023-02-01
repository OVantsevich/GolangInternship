package main

import (
	"context"
	"fmt"
	"github.com/OVantsevich/GolangInternship/FMicroservice/internal/config"
	"github.com/OVantsevich/GolangInternship/FMicroservice/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"testing"
)

var (
	cfg = config.Config{
		CurrentDB:   "postgres",
		PostgresUrl: "postgres://postgres:postgres@localhost:5432/userService?sslmode=disable",
		MongoURL:    "mongodb://mongo:mongo@localhost:27017",
		JwtKey:      "874967EC3EA3490F8F2EF6478B72A756",
	}
)

func DBConnectionTest() (repository.User, error) {
	switch cfg.CurrentDB {
	case "postgres":
		pool, err := pgxpool.New(context.Background(), cfg.PostgresUrl)
		if err != nil {
			return nil, fmt.Errorf("invalid configuration data: %v", err)
		}
		if err = pool.Ping(context.Background()); err != nil {
			return nil, fmt.Errorf("database not responding: %v", err)
		}
		return &repository.PUser{Pool: pool}, nil
	case "mongo":
		client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(cfg.MongoURL))
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

func ClosePoolTest(r interface{}) {
	switch cfg.CurrentDB {
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

func TestMain(m *testing.M) {
	log.Println("Do stuff BEFORE the tests!")
	exitVal := m.Run()
	log.Println("Do stuff AFTER the tests!")

	os.Exit(exitVal)
}
