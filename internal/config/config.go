package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type key string

const KeyMetrics = key("metrics")
const KeyUUID = key("uuid")

type Config struct {
	Service  Service
	Redis    Redis
	Platform Platform
	Metrics  Metrics
}

type Service struct {
	Port string `env:"CHAT_WORKER_PORT"`
}

type Redis struct {
	Host    string `env:"CHAT_REDIS_HOST"`
	Port    string `env:"CHAT_REDIS_PORT"`
	Channel string `env:"CHAT_REDIS_CHANNEL"`
}

//type Kafka struct {
//	NotificationNewFriendTopic string `env:"FRIENDS_EMAIL_INVITE"`
//	Server                     string `env:"KAFKA_SERVER"`
//	GroupID                    string `env:"KAFKA_GROUP_ID"`
//	AutoOffset                 string `env:"KAFKA_OFFSET"`
//}

type Metrics struct {
	Host string `env:"GRAFANA_HOST"`
	Port int    `env:"GRAFANA_PORT"`
}

type Platform struct {
	Env string `env:"ENV"`
}

func MustLoad() *Config {
	cfg := &Config{}
	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		log.Fatalf("failed to read env variables: %s", err)
	}
	return cfg
}
