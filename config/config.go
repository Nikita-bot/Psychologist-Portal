package config

import (
	"fmt"
	"log"

	"github.com/goloop/env"
)

type Config struct {
	Login    string `env:"LOGIN"`
	Password string `env:"PASS"`
	Port     string `env:"PORT"`
}

const (
	INDIVIDUAL = "Индивидуальная консультация"
	ROOM       = "Комната психологической разгрузки"
)

func InitConfig() *Config {
	if err := env.Load(".env"); err != nil {
		log.Fatal(err)
	}

	var cfg Config
	if err := env.Unmarshal("", &cfg); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", cfg)

	return &cfg
}
