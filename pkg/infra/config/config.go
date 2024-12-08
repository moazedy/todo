package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Postgres PostgresConfig
}

type ServerConfig struct {
	Address string
	Port    string
}

type PostgresConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Name     string
	Driver   string
}

func (pc PostgresConfig) ToString() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s port=%s sslmode=disable",
		pc.Host, pc.Username, pc.Password, pc.Port,
	)
}

func (pc PostgresConfig) ToStringWithDbName() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		pc.Host, pc.Username, pc.Password, pc.Name, pc.Port,
	)
}

func Init() Config {
	viper.SetConfigFile("./configs/vmanage.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		println(err.Error())
		panic("unable to read config file")
	}

	var config Config

	err = viper.Unmarshal(&config)
	if err != nil {
		println(err.Error())
		panic("unable to read config file")
	}

	return config
}
