package main

import "github.com/spf13/viper"

func readDefaultConfig() string {
	viper.SetDefault("port", ":3000")

	viper.SetConfigName("config")

	port := viper.Get("port").(string)

	return port
}