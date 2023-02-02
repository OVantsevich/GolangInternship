package config

type StreamConfig struct {
	StreamName string `env:"STREAM_NAME,notEmpty"`
}

func NewStreamConfig() (*StreamConfig, error) {
	Cfg := &StreamConfig{}
	//if err := env.Parse(Cfg); err != nil {
	//	return nil, fmt.Errorf("config - NewConfig: %v", err)
	//}
	Cfg.StreamName = "USERS" //😁
	return Cfg, nil
}
