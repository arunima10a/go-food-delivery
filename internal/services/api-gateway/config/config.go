package config

import (
	"log"
	"strings"
	"github.com/spf13/viper"
)

type Config struct {
	Service struct {
		Port string `mapstructure:"port"`
		Name string `mapstructure:"name"`
	} `mapstructure:"service"`
	Services struct {
		Catalog   string `mapstructure:"catalog"`
		Inventory string `mapstructure:"inventory"`
		Identity  string `mapstructure:"identity"`
		Ordering  string `mapstructure:"ordering"`
		Search    string `mapstructure:"search"`
	} `mapstructure:"services"`
}

func GetConfig() *Config {
	cfg := &Config{}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AddConfigPath(".../config")
	viper.AddConfigPath("./config")

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %v ,err")

	}

	if err := viper.Unmarshal(cfg); err != nil {
		log.Fatalf("Unable to decode into struct, %v, err")
	}
	return cfg
}
