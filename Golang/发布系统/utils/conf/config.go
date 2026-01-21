package config

import (
	"github.com/spf13/viper"
)

const (
	YAML  = "yaml"
	JSON  = "json"
	XML	  = "xml"
)

var Config *viper.Viper

func InitConfig(dir, file, fileType string) {
	config := viper.New()
	config.AddConfigPath(dir)
	config.SetConfigName(file)
	config.SetConfigType(fileType)

	if err := config.ReadInConfig(); err != nil {
		panic(err)
	}

	Config = config
}