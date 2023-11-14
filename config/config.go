package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	DbHost     string `mapstructure:"POSTGRES_HOST"`
	DbUser     string `mapstructure:"POSTGRES_USER"`
	DbPassword string `mapstructure:"POSTGRES_PASSWORD"`
	DbName     string `mapstructure:"POSTGRES_DB"`

	MqHost     string `mapstructure:"RABBITMQ_HOST"`
	MqUser     string `mapstructure:"RABBITMQ_DEFAULT_USER"`
	MqPassword string `mapstructure:"RABBITMQ_DEFAULT_PASS"`
}

func NewConfig() (Config, error) {
	config := Config{}
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	viper.AutomaticEnv()
	err = viper.Unmarshal(&config)
	if err != nil {
		return config, err
	}
	return config, nil
}
