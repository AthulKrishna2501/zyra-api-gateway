package main

import (
	"log"

	"github.com/AthulKrishna2501/zyra-api-gateway/pkg/config"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.LoadConfig()

	if err != nil {
		log.Fatal("Failed to config .env", err, cfg)
	}

	router := gin.Default()

	log.Print("Server start running on port:3000")
	router.Run(":3000")

}
