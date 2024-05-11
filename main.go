package main

import (
	"fmt"
	"kiit-lab-engine/routes"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	ginMode := viper.GetString("GIN_MODE")
	if ginMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	r := gin.Default()

	r.SetTrustedProxies(nil)
	routes.InitRoutes(r)

	port := viper.GetString("PORT")
	if port == "" {
		port = "8421"
	}

	if err := r.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Println(err.Error())
		return
	}

}
