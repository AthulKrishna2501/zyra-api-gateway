package clients

import (
	"log"

	pb "github.com/AthulKrishna2501/proto-repo/vendor"
	"github.com/AthulKrishna2501/zyra-api-gateway/internals/middleware"
	"github.com/AthulKrishna2501/zyra-api-gateway/internals/services"
	"github.com/AthulKrishna2501/zyra-api-gateway/pkg/config"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type VendorClient struct {
	Client pb.VendorSeviceClient
}

func InitVendorClient(c *config.Config) *VendorClient {
	conn, err := grpc.NewClient(c.VENDOR_SVC_URL, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultCallOptions(
		grpc.MaxCallRecvMsgSize(1024*1024*100),
		grpc.MaxCallSendMsgSize(1024*1024*100),
	))

	if err != nil {
		log.Fatal("Could not connect to vendor client", err)
	}

	return &VendorClient{
		Client: pb.NewVendorSeviceClient(conn),
	}

}

func RegisterVendorRoutes(eng *gin.Engine, cfg *config.Config) *VendorClient {
	vc := InitVendorClient(cfg)

	if vc.Client == nil {
		log.Fatal("Vendor Service Client is nil")
	}

	routes := eng.Group("/vendor")
	routes.Use(middleware.VendorAuthMiddleware(config.RedisClient))
	routes.POST("/request-category", vc.RequestCategory)
	routes.GET("/list-categories", vc.ListCategory)

	return vc
}

func (vc *VendorClient) RequestCategory(ctx *gin.Context) {
	services.RequestCategory(ctx, vc.Client)
}

func (vc *VendorClient) ListCategory(ctx *gin.Context) {
	services.ListCategory(ctx, vc.Client)
}
