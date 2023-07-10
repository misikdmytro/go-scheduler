package config

import "github.com/spf13/viper"

type ServerConfig struct {
	Host string
	Port string
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

type RabbitMQConfig struct {
	Host     string
	Port     string
	User     string
	Password string
}

type Config struct {
	Server   ServerConfig
	DB       DBConfig
	RabbitMQ RabbitMQConfig
}

func LoadConfig() (Config, error) {
	conf := viper.New()

	conf.SetConfigName("config")
	conf.SetConfigType("yml")
	conf.AddConfigPath(".")

	if err := conf.ReadInConfig(); err != nil {
		return Config{}, err
	}

	var config Config
	err := conf.Unmarshal(&config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}
