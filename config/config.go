package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type Mode string
type LogOutput string

const (
	Development Mode = "development"
	Production  Mode = "production"

	LogConsole LogOutput = "console"
	LogJson    LogOutput = "json"
	LogMulti   LogOutput = "all"
)

type Config struct {
	Server struct {
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
		Output LogOutput `yaml:"output"`
		Color  bool      `yaml:"color"`
		File   struct {
			Path string `yaml:"path"`
		} `yaml:"file"`
		TimeFormat string `yaml:"timeformat"`
		Encrypted  bool   `yaml:"encrypted"`
	} `yaml:"log"`
	DB struct {
		Drivers struct {
			Mysql struct {
				DbName string `yaml:"db_name"`
				Host   string `yaml:"host"`
				User   string `yaml:"user"`
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
