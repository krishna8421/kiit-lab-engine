package main

import (
	"kiit-lab-engine/core/server"
	"log"

	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}
}

func main() {
	if err := server.StartServer(); err != nil {
		panic(err)
	}
}
