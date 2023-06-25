package config

import (
	"log"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	GrpcPort      string `env:"GRPC_PORT" env-required:"true"`
	AccessAddress string `env:"ACCESS_ADDRESS" env-default:"localhost:50051"`
}

var config *Config
var onceConfig sync.Once

func GetConfig() *Config {
	onceConfig.Do(func() {
		err := godotenv.Load()
		if err != nil {
			log.Fatalln(err)
		}

		config = &Config{}

		err = cleanenv.ReadEnv(config)
		if err != nil {
			log.Fatalln(err)
		}
	})

	return config
}
