package config

import "github.com/spf13/viper"

type Config struct {
	REST_PORT string `mapstructure:"REST_PORT"`

	DB_HOST     string `mapstructure:"DB_HOST"`
	DB_PORT     string `mapstructure:"DB_PORT"`
	DB_USER     string `mapstructure:"DB_USER"`
	DB_PASSWORD string `mapstructure:"DB_PASSWORD"`
	DB_NAME     string `mapstructure:"DB_NAME"`
}

var ENV Config

func LoadConfig() error {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	viper.SetDefault("REST_PORT", "8080")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	viper.AutomaticEnv()

	// Unmarshal into the ENV variable
	if err := viper.Unmarshal(&ENV); err != nil {
		return err
	}

	return nil
}