package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	App struct {
		CertFilePath string `yaml:"cert_file_path"`
	}
	// For rest server
	Server struct {
		Host    string `yaml:"host"`
		Port    string `yaml:"port"`
		Timeout struct {
			Server time.Duration `yaml:"server"`
			Write  time.Duration `yaml:"write"`
			Read   time.Duration `yaml:"read"`
			Idle   time.Duration `yaml:"idle"`
		} `yaml:"timeout"`
	} `yaml:"server"`
	Database struct {
		Host           string `yaml:"host"`
		Port           string `yaml:"port"`
		User           string `yaml:"user"`
		Password       string `yaml:"password"`
		DBName         string `yaml:"dbname"`
		SSLMode        string `yaml:"sslmode"`
		PoolConnection struct {
			MaxOpenConns int           `yaml:"max_open_conns"`
			MaxIdleConns int           `yaml:"max_idle_conns"`
			MaxLifeTime  time.Duration `yaml:"max_life_time"`
		} `yaml:"pool_connection"`
	} `yaml:"database"`
	JWT struct {
		AccessTokenDuration  time.Duration `yaml:"accessTokenDuration"`
		RefreshTokenDuration time.Duration `yaml:"refreshTokenDuration"`
		Issuer               string        `yaml:"issuer"`
	} `yaml:"jwt"`
}

func NewConfig(file_name string) (*Config, error) {
	var config Config

	file, err := os.Open(file_name)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)

	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
