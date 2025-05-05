package healthcheck

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthCheckResponse struct {
	Status   string   `json:"status"`
	Services []string `json:"services"`
}

func checkGateway() string {
	return "API Gateway is up and running"
}

func HealthCheckHandler(c *gin.Context) {
	services := []string{
		checkGateway(),
	}
	response := HealthCheckResponse{
		Status:   "Healthy",
		Services: services,
	}

	c.JSON(http.StatusOK, response)
}
