package domain

type (
	Config struct {
		CurrentDB   string `env:"CURRENT_DB,notEmpty" envDefault:"postgres"`
		PostgresUrl string `env:"POSTGRES_DB_URL,notEmpty"`
		MongoURL    string `env:"MONGO_DB_URL,notEmpty"`
		JwtKey      string `env:"JWT_KEY,notEmpty"`
	}
)

var Cfg *Config

func InitConfig() error {
	if Cfg == nil {
		Cfg = &Config{}
		//if err := env.Parse(Cfg); err != nil {
		//	log.Fatalf("something went wrong with the environment: %v", err)
		//	return err
		//}

		Cfg.CurrentDB = "postgres"
		Cfg.PostgresUrl = "postgres://postgres:postgres@host.docker.internal:5432/entity?sslmode=disable"
		Cfg.MongoURL = "_"
		Cfg.JwtKey = "874967EC3EA3490F8F2EF6478B72A756"
		return nil
	}
	return nil
}
