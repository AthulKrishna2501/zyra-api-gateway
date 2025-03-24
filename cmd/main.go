package main

import (
	"log"

	"github.com/AthulKrishna2501/zyra-api-gateway/internals/clients"
	"github.com/AthulKrishna2501/zyra-api-gateway/pkg/config"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.LoadConfig()
	config.InitRedis()

	if err != nil {
		log.Fatal("Failed to config .env", err, cfg)
	}

	log.Print("Configurations loaded succesfully....")

	router := gin.Default()

	clients.RegisterAuthRoutes(router, &cfg)
	clients.RegisterVendorRoutes(router, &cfg)
	clients.RegisterAdminRoutes(router, &cfg)
	clients.RegisterClientClient(router, &cfg)

	log.Print("Server start running on port:3000")
	router.Run(":3000")

}
