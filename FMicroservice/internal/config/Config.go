package config

import (
	"fmt"
	"github.com/caarlos0/env/v6"
)

type (
	Config struct {
		CurrentDB   string `env:"CURRENT_DB,notEmpty" envDefault:"postgres"`
		PostgresUrl string `env:"POSTGRES_DB_URL,notEmpty"`
		MongoURL    string `env:"MONGO_DB_URL,notEmpty"`
		JwtKey      string `env:"JWT_KEY,notEmpty"`
	}
)

func NewConfig() (*Config, error) {
	Cfg := &Config{}
	if err := env.Parse(Cfg); err != nil {
		return nil, fmt.Errorf("config - NewConfig: %v", err)
	}

	//Cfg.CurrentDB = "mongo"
	//Cfg.PostgresUrl = "postgres://postgres:postgres@host.docker.internal:5432/entity?sslmode=disable"
	//Cfg.MongoURL = "mongodb://mongo:mongo@host.docker.internal:27017"
	//Cfg.JwtKey = "874967EC3EA3490F8F2EF6478B72A756"
	return Cfg, nil
}
