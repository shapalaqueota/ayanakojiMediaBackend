package config

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	BitPayAPIKey string `mapstructure:"BITPAY_API_KEY"`
}

var AppConfig Config

func LoadConfig() {
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	if err := viper.Unmarshal(&AppConfig); err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}
}
