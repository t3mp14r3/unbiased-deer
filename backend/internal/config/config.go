package config

import (
	"log"

	env "github.com/caarlos0/env/v6"
)

type Config struct {
    PostgresConfig      PostgresConfig
    RedisConfig         RedisConfig
    NatsConfig          NatsConfig
    GatewayConfig       GatewayConfig
}

type PostgresConfig struct {
    Host        string `env:"POSTGRES_HOST"`
    Port        int    `env:"POSTGRES_PORT"`
    Name        string `env:"POSTGRES_NAME"`
    User        string `env:"POSTGRES_USER"`
    Password    string `env:"POSTGRES_PASS"`
}

type RedisConfig struct {
    Addr        string `env:"REDIS_ADDR"`
    Password    string `env:"REDIS_PASS"`
    DB          int    `env:"REDIS_DB"`
}

type NatsConfig struct {
    Url         string `env:"NATS_URL"`
}

type GatewayConfig struct {
    Addr        string `env:"GATEWAY_ADDR"`
}

func New() *Config {
    var config Config

	if err := env.Parse(&config.PostgresConfig); err != nil {
        log.Fatalln("failed to parse postgres config! err:", err)
	}
	
    if err := env.Parse(&config.RedisConfig); err != nil {
        log.Fatalln("failed to parse redis config! err:", err)
	}
	
    if err := env.Parse(&config.NatsConfig); err != nil {
        log.Fatalln("failed to parse nats config! err:", err)
	}
	
    if err := env.Parse(&config.GatewayConfig); err != nil {
        log.Fatalln("failed to parse gateway server config! err:", err)
	}

	return &config
}

