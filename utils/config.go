package utils

import (
	"github.com/hashicorp/go-hclog"
	"github.com/spf13/viper"
)

type Config struct {
	Environment string

	Server struct {
		Host string
		Grpc struct {
			Port string
		}
		Gateway struct {
			Port string
		}
	}

	Database struct {
		Name       string
		URI        string
		Collection struct {
			Video string
		}
	}

	Aws struct {
		Video struct {
			Bucket string
		}
		Thumbnail struct {
			Bucket string
		}
		Region string
	}
}

var configuration *Config = nil

func LoadConfig(l hclog.Logger) (*Config, error) {
	var config *Config

	viper.AddConfigPath("config/")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		l.Error("Unable to read the config file: ", err.Error())
		return nil, err
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		l.Error("Unable to unmarshall: ", err.Error())
		return nil, err
	}
	configuration = config

	return config, nil
}

func GetConfig() *Config {
	return configuration
}
