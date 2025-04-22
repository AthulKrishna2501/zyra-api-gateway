package clients

import (
	"log"
	"time"

	pb "github.com/AthulKrishna2501/proto-repo/client"
	"github.com/AthulKrishna2501/zyra-api-gateway/internals/middleware"
	"github.com/AthulKrishna2501/zyra-api-gateway/internals/services"
	"github.com/AthulKrishna2501/zyra-api-gateway/pkg/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sony/gobreaker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ClientClient struct {
	Client pb.ClientServiceClient
	Cfg    *config.Config
	CB     *gobreaker.CircuitBreaker
}

func newCircuitBreaker() *gobreaker.CircuitBreaker {
	return gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        "ClientServiceCB",
		MaxRequests: 5,
		Interval:    10 * time.Second,
		Timeout:     5 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures > 3
		},
	})
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
		CB:     newCircuitBreaker(),
		Cfg:    c,
	}

}

func RegisterClientClient(eng *gin.Engine, cfg *config.Config) *ClientClient {
	cc := InitClientClient(cfg)

	if cc.Client == nil {
		log.Fatal("Client Service Client is nil")
	}

	eng.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3005"},
		AllowMethods: []string{"GET", "POST"},
		AllowHeaders: []string{"Origin", "Content-Type"},
	}))

	routes := eng.Group("/client")
	routes.Use(middleware.ClientAuthMiddleware(config.RedisClient))
	routes.POST("/mc/payment", cc.PayMasterOfCeremony)
	routes.POST("/host-event", cc.HostEvent)
	routes.PUT("/edit-event", cc.EditEvent)
	routes.GET("/profile", cc.ClientProfile)
	routes.PUT("/profile", cc.EditClientProfile)
	routes.PUT("/reset-password", cc.ResetPassword)
	routes.GET("/bookings", cc.GetBookings)
	routes.GET("/dashboard", cc.ClientDashboard)
	routes.POST("/booking", cc.BookVendor)
	routes.GET("/vendors", cc.GetVendorsByCategory)
	routes.GET("/hosted-events", cc.GetHostedEvents)
	routes.GET("/upcoming-events", cc.GetUpcomingEvents)
	routes.GET("/vendor-profile", cc.GetVendorProfile)

	eng.POST("/webhook", cc.HandleStripeWebhook)

	return cc
}

func (cc *ClientClient) PayMasterOfCeremony(ctx *gin.Context) {
	state := cc.CB.State().String()
	log.Println("Circuit Breaker State (Before Call):", state)

	_, err := cc.CB.Execute(func() (interface{}, error) {
		services.PayMasterOfCeremony(ctx, cc.Client)
		return nil, nil
	})

	state = cc.CB.State().String()
	log.Println("Circuit Breaker State (After Call):", state)

	if err != nil {
		ctx.JSON(503, gin.H{"error": "Payment Service Unavailable"})
		return
	}
}

func (cc *ClientClient) HandleStripeWebhook(ctx *gin.Context) {
	_, err := cc.CB.Execute(func() (interface{}, error) {
		services.HandleStripeWebhook(ctx, cc.Client, cc.Cfg)
		return nil, nil
	})

	if err != nil {
		ctx.JSON(503, gin.H{"error": "Payment Service Unavailable"})
		return
	}
}


func (cc *ClientClient) HostEvent(ctx *gin.Context) {
	_, err := cc.CB.Execute(func() (interface{}, error) {
		services.HostEvent(ctx, cc.Client)
		return nil, nil

	})

	if err != nil {
		ctx.JSON(503, gin.H{"error": "Client Service Unavailable"})
		return
	}
}

func (cc *ClientClient) EditEvent(ctx *gin.Context) {
	_, err := cc.CB.Execute(func() (interface{}, error) {
		services.EditEvent(ctx, cc.Client)
		return nil, nil

	})

	if err != nil {
		ctx.JSON(503, gin.H{"error": "Client Service Unavailable"})
		return

	}
}

func (cc *ClientClient) ClientProfile(ctx *gin.Context) {
	_, err := cc.CB.Execute(func() (interface{}, error) {
		services.GetClientProfile(ctx, cc.Client)
		return nil, nil

	})

	if err != nil {
		ctx.JSON(503, gin.H{"error": "Client Service Unavailable"})
		return

	}
}

func (cc *ClientClient) EditClientProfile(ctx *gin.Context) {
	_, err := cc.CB.Execute(func() (interface{}, error) {
		services.EditClientProfile(ctx, cc.Client)
		return nil, nil

	})

	if err != nil {
		ctx.JSON(503, gin.H{"error": "Client Service Unavailable"})
		return

	}
}

func (cc *ClientClient) ResetPassword(ctx *gin.Context) {
	_, err := cc.CB.Execute(func() (interface{}, error) {
		services.ResetPassword(ctx, cc.Client)
		return nil, nil

	})

	if err != nil {
		ctx.JSON(503, gin.H{"error": "Client Service Unavailable"})
		return

	}
}

func (cc *ClientClient) GetBookings(ctx *gin.Context) {
	_, err := cc.CB.Execute(func() (interface{}, error) {
		services.GetBookings(ctx, cc.Client)
		return nil, nil

	})

	if err != nil {
		ctx.JSON(503, gin.H{"error": "Client Service Unavailable"})
		return

	}
}

func (cc *ClientClient) ClientDashboard(ctx *gin.Context) {
	_, err := cc.CB.Execute(func() (interface{}, error) {
		services.ClientDashboard(ctx, cc.Client)
		return nil, nil

	})

	if err != nil {
		ctx.JSON(503, gin.H{"error": "Client Service Unavailable"})
		return

	}
}

func (cc *ClientClient) BookVendor(ctx *gin.Context) {
	_, err := cc.CB.Execute(func() (interface{}, error) {
		services.BookVendor(ctx, cc.Client)
		return nil, nil

	})

	if err != nil {
		ctx.JSON(503, gin.H{"error": "Client Service Unavailable"})
		return

	}
}

func (cc *ClientClient) GetVendorsByCategory(ctx *gin.Context) {
	_, err := cc.CB.Execute(func() (interface{}, error) {
		services.GetVendorsByCategory(ctx, cc.Client)
		return nil, nil

	})

	if err != nil {
		ctx.JSON(503, gin.H{"error": "Client Service Unavailable"})
		return

	}
}

func (cc *ClientClient) GetHostedEvents(ctx *gin.Context) {
	_, err := cc.CB.Execute(func() (interface{}, error) {
		services.GetHostedEvents(ctx, cc.Client)
		return nil, nil

	})

	if err != nil {
		ctx.JSON(503, gin.H{"error": "Client Service Unavailable"})
		return

	}
}

func (cc *ClientClient) GetUpcomingEvents(ctx *gin.Context) {
	_, err := cc.CB.Execute(func() (interface{}, error) {
		services.GetUpcomingEvents(ctx, cc.Client)
		return nil, nil

	})

	if err != nil {
		ctx.JSON(503, gin.H{"error": "Client Service Unavailable"})
		return

	}
}

func (cc *ClientClient) GetVendorProfile(ctx *gin.Context) {
	_, err := cc.CB.Execute(func() (interface{}, error) {
		services.GetVendorProfile(ctx, cc.Client)
		return nil, nil

	})

	if err != nil {
		ctx.JSON(503, gin.H{"error": "Client Service Unavailable"})
		return

	}
}
