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

	// Authentication
	routes := eng.Group("/auth")
	routes.POST("/register", svc.Register)
	routes.POST("/send-otp", svc.SendOTP)
	routes.POST("/login", svc.Login)
	routes.GET("/google-login", svc.GoogleLogin)
	routes.GET("/callback", svc.HandleGoogleCallback)
	routes.POST("/verify-otp", svc.VerifyOTP)
	routes.POST("/resend-otp", svc.ResendOTP)
	routes.GET("/refresh-token", svc.RefreshToken)
	routes.POST("/logout", svc.Logout)

	return svc
}

func (svc *ServiceClient) Register(ctx *gin.Context) {
	services.Register(ctx, svc.Client)
}

func (svc *ServiceClient) SendOTP(ctx *gin.Context) {
	services.SendOTP(ctx, svc.Client)
}

func (svc *ServiceClient) VerifyOTP(ctx *gin.Context) {
	services.VerifyOTP(ctx, svc.Client)
}

func (svc *ServiceClient) ResendOTP(ctx *gin.Context) {
	services.ResendOTP(ctx, svc.Client)
}

func (svc *ServiceClient) Login(ctx *gin.Context) {
	services.Login(ctx, svc.Client)
}

func (svc *ServiceClient) RefreshToken(ctx *gin.Context) {
	services.RefreshToken(ctx, svc.Client)
}

func (svc *ServiceClient) Logout(ctx *gin.Context) {
	services.Logout(ctx, svc.Client)
}

func (svc *ServiceClient) GoogleLogin(ctx *gin.Context) {
	services.GoogleLogin(ctx, svc.Client)
}

func (svc *ServiceClient) HandleGoogleCallback(ctx *gin.Context) {
	services.HandleGoogleCallback(ctx, svc.Client)
}
