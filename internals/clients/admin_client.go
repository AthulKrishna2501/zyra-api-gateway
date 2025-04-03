package clients

import (
	"log"

	pb "github.com/AthulKrishna2501/proto-repo/admin"
	"github.com/AthulKrishna2501/zyra-api-gateway/internals/middleware"
	"github.com/AthulKrishna2501/zyra-api-gateway/internals/services"
	"github.com/AthulKrishna2501/zyra-api-gateway/pkg/config"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AdminClient struct {
	Client pb.AdminServiceClient
}

func InitAdminClient(c *config.Config) *AdminClient {
	conn, err := grpc.NewClient(c.ADMIN_SVC_URL, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultCallOptions(
		grpc.MaxCallRecvMsgSize(1024*1024*100),
		grpc.MaxCallSendMsgSize(1024*1024*100),
	))

	if err != nil {
		log.Fatal("Could not connect to admin client", err)
	}

	return &AdminClient{
		Client: pb.NewAdminServiceClient(conn),
	}
}

func RegisterAdminRoutes(eng *gin.Engine, cfg *config.Config) *AdminClient {
	ac := InitAdminClient(cfg)
	if ac.Client == nil {
		log.Fatal("Admin Service Client is nil")
	}

	routes := eng.Group("/admin")
	routes.Use(middleware.AdminAuthMiddleware(config.RedisClient))
	routes.POST("/approve-reject", ac.ApproveRejectCategory)
	routes.PUT("/block-user", ac.BlockUser)
	routes.PUT("/unblock-user", ac.UnblockUser)
	routes.GET("/users", ac.ListUsers)
	routes.GET("/view-requests", ac.ViewCategoryRequests)
	routes.POST("/add-category", ac.AddCategory)
	routes.GET("/dashboard", ac.AdminDashboard)
	routes.GET("/wallet", ac.GetAdminWallet)

	return ac
}

func (ac *AdminClient) ApproveRejectCategory(ctx *gin.Context) {
	services.ApproveRejectCategory(ctx, ac.Client)
}

func (ac *AdminClient) BlockUser(ctx *gin.Context) {
	services.BlockUser(ctx, ac.Client)
}

func (ac *AdminClient) UnblockUser(ctx *gin.Context) {
	services.UnblockUser(ctx, ac.Client)
}

func (ac *AdminClient) ListUsers(ctx *gin.Context) {
	services.ListUsers(ctx, ac.Client)
}

func (ac *AdminClient) ViewCategoryRequests(ctx *gin.Context) {
	services.ViewRequests(ctx, ac.Client)
}

func (ac *AdminClient) AddCategory(ctx *gin.Context) {
	services.AddCategory(ctx, ac.Client)
}

func (ac *AdminClient) AdminDashboard(ctx *gin.Context) {
	services.AdminDashboard(ctx, ac.Client)
}

func (ac *AdminClient) GetAdminWallet(ctx *gin.Context) {
	services.GetAdminWallet(ctx, ac.Client)
}
