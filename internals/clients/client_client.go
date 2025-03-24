package clients

import (
	"log"

	pb "github.com/AthulKrishna2501/proto-repo/client"
	"github.com/AthulKrishna2501/zyra-api-gateway/internals/middleware"
	"github.com/AthulKrishna2501/zyra-api-gateway/internals/services"
	"github.com/AthulKrishna2501/zyra-api-gateway/pkg/config"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ClientClient struct {
	Client pb.ClientServiceClient
	Cfg    config.Config
}

func InitClientClient(c *config.Config) *ClientClient {
	conn, err := grpc.NewClient(c.CLIENT_SVC_URL, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultCallOptions(
		grpc.MaxCallRecvMsgSize(1024*1024*100),
		grpc.MaxCallSendMsgSize(1024*1024*100),
	))

	if err != nil {
		log.Fatal("Could not connect to client client", err)
	}

	return &ClientClient{
		Client: pb.NewClientServiceClient(conn),
	}

}

func RegisterClientClient(eng *gin.Engine, cfg *config.Config) *ClientClient {
	cc := InitClientClient(cfg)

	if cc.Client == nil {
		log.Fatal("Client Service Client is nil")
	}

	routes := eng.Group("/client")
	routes.Use(middleware.ClientAuthMiddleware(config.RedisClient))
	routes.POST("/mc/payment", cc.PayMasterOfCeremony)
	routes.POST("/webhook", cc.HandleStripeWebhook)

	return cc
}

func (cc *ClientClient) PayMasterOfCeremony(ctx *gin.Context) {
	services.PayMasterOfCeremony(ctx, cc.Client)
}

func (cc *ClientClient) HandleStripeWebhook(ctx *gin.Context) {
	services.HandleStripeWebhook(ctx, cc.Client, cc.Cfg)
}
