package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"path/filepath"
	"time"
)

type Config struct {
	Env            string `yaml:"env" env-default:"development"`
	SecretKey      string `yaml:"secret_key" env-required:"true"`
	HTTPServer     `yaml:"http_server"`
	DatabaseConfig `yaml:"database"`
}

type DatabaseConfig struct {
	Name     string `yaml:"name"     env-required:"true"`
	User     string `yaml:"user"     env-required:"true"`
	Password string `yaml:"password" env-required:"true"`
	Port     string `yaml:"port"     env-required:"true"`
	Host     string `yaml:"host"     env-required:"true"`
}

type HTTPServer struct {
	Address     string        `yaml:"address"      env-default:"0.0.0.0:8080"`
	Timeout     time.Duration `yaml:"timeout"      env-default:"5s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func MustLoad() *Config {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatalf("error getting executable path: %s", err)
	}

	exeDir := filepath.Dir(exePath)

	configPath := filepath.Join(exeDir, "config.yml")
	fmt.Println("PATH ->", configPath)
	if _, err := os.Stat(configPath); err != nil {
		log.Fatalf("error opening config file: %s", err)
	}

	var cfg Config

	err = cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("error reading config file: %s", err)
	}

	return &cfg
}
