package config

import (
	"log"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Logger     LoggerConf
	DB         DBConf
	Storage    StorageConf
	HTTPServer ServerConf
	GrpcServer ServerConf
	RabbitMq   RabbitMqConf
}

type LoggerConf struct {
	Level string
}

type DBConf struct {
	Dsn string
}

type RabbitMqConf struct {
	URL string
}

type StorageConf struct {
	Type string
}

type ServerConf struct {
	Host string
	Port string
}

func NewConfig(configFilePath string) Config {
	var config Config
	_, err := toml.DecodeFile(configFilePath, &config)
	if err != nil {
		log.Fatalln("Error decode config file...")
	}
	return config
}
