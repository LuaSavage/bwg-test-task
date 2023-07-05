package config

import (
	"sync"

	"github.com/LuaSavage/bwg-test-task/service-b/pkg/logging"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	IsDebug *bool `yaml:"is_debug" env-required:"true"`
	Listen  struct {
		Type   string `yaml:"type" env-default:"port"`
		BindIP string `yaml:"bind_ip" env-default:"127.0.0.1"`
		Port   string `yaml:"port" env-default:"8080"`
	} `yaml:"listen"`
	Storage StorageConfig `yaml:"storage"`
	Kafka   KafkaConfig   `yaml:"kafka"`
}

type StorageConfig struct {
	Host     string `json:"host" env:"PSQL_HOST"`
	Port     string `json:"port" env:"PSQL_PORT"`
	Database string `json:"database" env:"PSQL_DATABASE"`
	Username string `json:"username" env:"PSQL_USERNAME"`
	Password string `json:"password" env:"PSQL_PASSWORD"`
}

type KafkaConfig struct {
	Host  string `yaml:"host" env:"KAFKA_HOST" env-default:"localhost"`
	Port  string `yaml:"port" env:"KAFKA_PORT" env-default:"9092"`
	Topic string `yaml:"password" env:"KAFKA_TOPIC" env-required:"true" env-default:"transfer"`
}

var instance *Config
var once sync.Once

func GetConfig(path string) (*Config, error) {
	var err error
	once.Do(func() {
		logger := logging.GetLogger()
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
