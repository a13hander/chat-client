package config

import (
	"log"
	"sync"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	AuthAddr             string        `env:"AUTH_ADDRESS" env-default:"localhost:50051"`
	ChatSeverAddr        string        `env:"CHAT_ADDRESS" env-default:"localhost:50052"`
	RedisAddr            string        `env:"REDIS_ADDRESS" env-default:"localhost:6379"`
	RefreshTokenLifeTime time.Duration `env:"REFRESH_TOKEN_LIFETIME" env-default:"1h"`
	AccessTokenLifeTime  time.Duration `env:"ACCESS_TOKEN_LIFETIME" env-default:"1m"`
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
