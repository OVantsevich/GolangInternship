// Package config main config
package config

// Config with init data
type Config struct {
	CurrentDB   string `env:"CURRENT_DB,notEmpty" envDefault:"postgres"`
	PostgresURL string `env:"POSTGRES_DB_URL,notEmpty"`
	MongoURL    string `env:"MONGO_DB_URL,notEmpty"`
	JwtKey      string `env:"JWT_KEY,notEmpty"`
}

// NewConfig Creating Config using env parsing
func NewConfig() (*Config, error) {
	Cfg := &Config{}

	Cfg.CurrentDB = "postgres"
	Cfg.PostgresURL = "postgres://postgres:postgres@localhost:5432/userService?sslmode=disable"
	Cfg.MongoURL = "mongodb://mongo:mongo@localhost:27017"
	Cfg.JwtKey = "874967EC3EA3490F8F2EF6478B72A756"
	return Cfg, nil
}
