package config

import (
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/labstack/echo/v4"
)

type Config struct {
	Listen struct {
		BindIP string `yaml:"bind_ip" env:"BIND_IP" env-default:"127.0.0.1"`
		Port   string `yaml:"port" env:"PORT" env-default:"8080"`
		RPC    int    `yaml:"rpc" env:"RATE_LIMIT" env-default:"10"`
	} `yaml:"listen"`
	Kafka struct {
		Host  string `yaml:"host" env:"KAFKA_HOST" env-default:"localhost"`
		Port  string `yaml:"port" env:"KAFKA_PORT" env-default:"9092"`
		Topic string `yaml:"topic" env:"KAFKA_TOPIC" env-required:"true" env-default:"transfer"`
	} `yaml:"kafka"`
}

var (
	instance *Config
	once     sync.Once
)

func GetConfig(path string, logger echo.Logger) (*Config, error) {
	var err error

	once.Do(func() {
		logger.Info("read application config")
		instance = &Config{}

		if len(path) == 0 {
			path = "../config.yml"
		}

		if err = cleanenv.ReadConfig(path, instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Info(help)
		}
	})

	return instance, err
}
