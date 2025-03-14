package clients

import (
	"log"

	pb "github.com/AthulKrishna2501/proto-repo/auth"
	"github.com/AthulKrishna2501/zyra-api-gateway/internals/services"
	"github.com/AthulKrishna2501/zyra-api-gateway/pkg/config"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ServiceClient struct {
	Client pb.AuthServiceClient
}

func InitServiceClient(c *config.Config) *ServiceClient {
	conn, err := grpc.NewClient(c.AUTH_SVC_URL, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatal("Could not connect to auth client", err)
	}

	return &ServiceClient{
		Client: pb.NewAuthServiceClient(conn),
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
	routes.POST("/verify-otp", svc.VerifyOTP)
	routes.POST("/logout", svc.Logout)

	return svc
}

func (svc *ServiceClient) Register(ctx *gin.Context) {
	services.Register(ctx, svc.Client)
}

func (svc *ServiceClient) Login(ctx *gin.Context) {
	services.Login(ctx, svc.Client)
}

func (svc *ServiceClient) VerifyOTP(ctx *gin.Context) {
	services.VerifyOTP(ctx, svc.Client)
}

func (svc *ServiceClient) Logout(ctx *gin.Context) {
	services.Logout(ctx, svc.Client)
}
