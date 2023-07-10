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

type JobsConfig struct {
	Exchange string
}

type Config struct {
	Server   ServerConfig
	DB       DBConfig
	RabbitMQ RabbitMQConfig
	Jobs     JobsConfig
}

func LoadConfig() (Config, error) {
	return LoadConfigFrom(".")
}

func LoadConfigFrom(in string) (Config, error) {
	conf := viper.New()

	conf.SetConfigName("config")
	conf.SetConfigType("yml")
	conf.AddConfigPath(in)

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
