package config

import "github.com/spf13/viper"

type RabbitMQConfig struct {
	Host     string
	Port     string
	User     string
	Password string
}

type PublisherConfig struct {
	Exchange string
}

type SubscriberConfig struct {
	Group string
	Queue string
}

type Config struct {
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
