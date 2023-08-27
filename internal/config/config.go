package config

import (
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	HTTPServer `yaml:"http_server"`
	DB         `yaml:"db"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:":8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"1s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"120s"`
}

type DB struct {
	Host         string `yaml:"host" env-default:"postgres"`
	DBName       string `yaml:"dbname" env-required:"true"`
	Username     string `yaml:"username" env-required:"true"`
	Password     string `env:"DB_PASSWORD"`
	Port         string `yaml:"port" env-default:"5432"`
	SSLMode      string `yaml:"sslmode" env-default:"disable"`
	InitFilepath string `yaml:"init_filepath" env-required:"true"`
}

func LoadConfig(configPath string) (*Config, error) {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file is not found in the specified path: %s", configPath)
	}
	config := new(Config)
	if err := cleanenv.ReadConfig(configPath, config); err != nil {
		return nil, fmt.Errorf("cannot load config: %s", err)
	}
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("cannot load env vars: %s", err)
	}
	if err := cleanenv.ReadEnv(config); err != nil {
		return nil, fmt.Errorf("cannot load config: %s", err)
	}
	return config, nil
}

func (d DB) String() string {
	return fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		d.Host, d.Port, d.Username, d.DBName, d.Password, d.SSLMode)
}
