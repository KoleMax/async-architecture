package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

var config *Config

type Config struct {
	Logging struct {
		Level string `yaml:"level"`
	} `yaml:"logging"`
	Db struct {
		User        string `yaml:"user"`
		Password    string `yaml:"password"`
		Host        string `yaml:"host"`
		Port        uint16 `yaml:"port"`
		Name        string `yaml:"name"`
		MinConnsNum int    `yaml:"min_conns_num" split_words:"true"`
		MaxConnsNum int    `yaml:"max_conns_num" split_words:"true"`
	} `yaml:"db"`
	Storage struct {
		Endpoint  string `yaml:"endpoint"`
		SecretKey string `yaml:"secret_key" split_words:"true"`
		AccessKey string `yaml:"access_key" split_words:"true"`
		Bucket    string `yaml:"user"`
	} `yaml:"storage"`
	Auth struct {
		IssuerUrl  string `yaml:"issuer_url" split_words:"true"`
		ClientName string `yaml:"client_name" split_words:"true"`
	} `yaml:"auth"`
}

func Get() *Config {
	return config
}

func init() {
	newConfigObj, err := createConfig()
	if err != nil {
		log.Fatal(err)
	}
	config = newConfigObj
}

func createConfig() (*Config, error) {
	cfg := &Config{}

	err := fillFromEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func fillFromEnv(config *Config) error {
	return envconfig.Process("", config)
}
