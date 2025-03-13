package clients

import (
	"log"

	pb "github.com/AthulKrishna2501/proto-repo/auth"
	"github.com/AthulKrishna2501/zyra-api-gateway/internals/events"
	"github.com/AthulKrishna2501/zyra-api-gateway/internals/services"
	"github.com/AthulKrishna2501/zyra-api-gateway/pkg/config"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ServiceClient struct {
	Client   pb.AuthServiceClient
	RabbitMq *events.RabbitMq
}

func InitServiceClient(c *config.Config) *ServiceClient {
	conn, err := grpc.NewClient(c.AUTH_SVC_URL, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatal("Could not connect to auth client", err)
	}

	rabbitMQ, err := events.NewRabbitMq(c.RABBITMQ_URL)
	if err != nil {
		log.Fatal("Could not connect to RabbitMQ:", err)
	}

	return &ServiceClient{
		Client:   pb.NewAuthServiceClient(conn),
		RabbitMq: rabbitMQ,
	}
}

func RegisterAuthRoutes(eng *gin.Engine, cfg *config.Config) *ServiceClient {
	svc := InitServiceClient(cfg)
	if svc.Client == nil {
		log.Fatal("Auth Service Client is nil!")
	}

	routes := eng.Group("/auth")
	routes.POST("/register", svc.Register)
	routes.POST("/login", svc.Login)
	routes.POST("/send-otp", svc.SendOTP)
	routes.POST("/verify-otp", svc.VerifyOTP)
	routes.POST("/logout", svc.Logout)

	return svc
}

func (svc *ServiceClient) Register(ctx *gin.Context) {
	services.Register(ctx, svc.Client, svc.RabbitMq)
}

func (svc *ServiceClient) Login(ctx *gin.Context) {
	services.Login(ctx, svc.Client)
}

func (svc *ServiceClient) SendOTP(ctx *gin.Context) {
	services.SendOTP(ctx, svc.Client)
}
func (svc *ServiceClient) VerifyOTP(ctx *gin.Context) {
	services.VerifyOTP(ctx, svc.Client)
}

func (svc *ServiceClient) Logout(ctx *gin.Context) {
	services.Logout(ctx, svc.Client)
}
