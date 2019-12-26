package config

import (
	"fmt"
	"github.com/jinzhu/configor"
)

type (
	LogConfig struct {
		//日志
		Path  string `yaml:"path" env:"LOG_PATH"`
		File  string `yaml:"file" env:"LOG_FILE"`
		Level int    `yaml:"level" env:"LOG_LEVEL"`
	}

	DbConfig struct {
		//数据库
		Host     string `yaml:"host" env:"DB_HOST"`
		Port     int    `yaml:"port" env:"DB_PORT"`
		User     string `yaml:"user" env:"DB_USER"`
		Password string `yaml:"password" env:"DB_PASSWORD"`
		Dbname   string `yaml:"dbname" env:"DB_DBNAME"`
	}

	Eth struct {
		Node string `yaml:"node" env:"NODE"`
	}

	Configuration struct {
		Log LogConfig `yaml:"log"`
		Db  DbConfig  `yaml:"db"`
		Eth Eth       `yaml:"eth"`
	}
)

var Cfg Configuration

func Init(filePath string) error {
	fmt.Println(filePath)
	err := configor.Load(&Cfg, "D:/etherscan-go/build/config.yaml")
	fmt.Printf("config: %#v\n", Cfg)
	return err
}
