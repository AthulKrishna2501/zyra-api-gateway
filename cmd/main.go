package main

import (
	"log"

	"github.com/AthulKrishna2501/zyra-api-gateway/internals/clients"
	"github.com/AthulKrishna2501/zyra-api-gateway/pkg/config"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.LoadConfig()

	if err != nil {
		log.Fatal("Failed to config .env", err, cfg)
	}

	router := gin.Default()

	auth:=clients.RegisterAuthRoutes(router,&cfg)

	log.Print(auth.Client)



	log.Print("Server start running on port:3000")
	router.Run(":3000")

}
