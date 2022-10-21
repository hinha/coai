package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type Mode string

type Config struct {
	Server struct {
		Grpc struct {
			Host string `yaml:"host"`
			Port int    `yaml:"port"`
		} `yaml:"grpc"`
		Host      string `yaml:"host"`
		Name      string `yaml:"name"`
		Port      int    `yaml:"port"`
		Timeout   int    `yaml:"timeout"`
		Mode      Mode   `yaml:"mode"`
		SecretKey string `yaml:"secret_key"`
	} `yaml:"server"`
	Profiler struct {
		Enabled bool   `yaml:"enabled"`
		Server  string `yaml:"server"`
	} `yaml:"profiler"`
	Jwt Jwt `yaml:"jwt"`
	Log struct {
		Color     bool `yaml:"color"`
		Encrypted bool `yaml:"encrypted"`
		File      struct {
			Http string `yaml:"path"`
			Grpc string `yaml:"grpc"`
		} `yaml:"file"`
		Output     string `yaml:"output"`
		TimeFormat string `yaml:"timeformat"`
	} `yaml:"log"`
	DB struct {
		Drivers struct {
			Mysql struct {
				DbName string `yaml:"db_name"`
				Host   string `yaml:"host"`
				User   string `yaml:"username"`
				Pass   string `yaml:"password"`
				Port   int    `yaml:"port"`
			} `yaml:"mysql"`
		} `yaml:"drivers"`
	} `yaml:"database"`
	Otel struct {
		Enabled bool   `yaml:"enabled"`
		Server  string `yaml:"server"`
	}
}

type Jwt struct {
	Expire     int    `yaml:"expire"`
	PrivateKey string `json:"private_key" yaml:"private_key"`
	PublicKey  string `json:"public_key" yaml:"public_key"`
	SecretKey  string `json:"secret_key" yaml:"secret_key"`
}

// LoadSecret reads the file from path and return Secret
func LoadSecret() (*Config, error) {
	data, err := ioutil.ReadFile("./config/config.yml")
	if err != nil {
		return nil, err
	}
	return LoadSecretFromBytes(data)
}

// LoadSecretFromBytes reads the secret file from data bytes
func LoadSecretFromBytes(data []byte) (*Config, error) {
	var cfg Config
	err := yaml.Unmarshal([]byte(data), &cfg)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return &cfg, nil
}
