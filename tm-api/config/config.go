package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Debug            bool   `envconfig:"debug"`
	Port             int    `envconfig:"port"`
	PostgresHost     string `envconfig:"postgres_host"`
	PostgresUser     string `envconfig:"postgres_user"`
	PostgresPassword string `envconfig:"postgres_password"`
	PostgresPort     int    `envconfig:"postgres_port"`
	PostgresDB       string `envconfig:"postgres_db"`
	Env              string `envconfig:"env"`
	Host             string `envconfig:"host"`
}

func Load() (*Config, error) {
	env := os.Getenv("GIN_MODE")
	if env != "release" {
		if err := godotenv.Load("./.env"); err != nil {
			log.Printf("couldn't load env vars: %v", err)
		}
	}

	c := &Config{}
	err := envconfig.Process("orchestra", c)
	if err != nil {
		return nil, err
	}
	return c, nil
}
