package main

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
}

type LoggerConf struct {
	Level string
}

type DBConf struct {
	Dsn string
}

type StorageConf struct {
	Type string
}

type ServerConf struct {
	Host string
	Port string
}

func NewConfig() Config {
	var config Config
	_, err := toml.DecodeFile(configFile, &config)
	if err != nil {
		log.Fatalln("Error decode config file...")
	}
	return config
}
