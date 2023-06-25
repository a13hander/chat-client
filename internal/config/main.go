package config

import (
	"log"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	AuthAddr      string `env:"ACCESS_ADDRESS" env-default:"localhost:50051"`
	ChatSeverAddr string `env:"ACCESS_ADDRESS" env-default:"localhost:50052"`
	RedisAddr     string `env:"ACCESS_ADDRESS" env-default:"localhost:6379"`
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
