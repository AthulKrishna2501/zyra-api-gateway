package services

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"

	pb "github.com/AthulKrishna2501/proto-repo/client"
	"github.com/AthulKrishna2501/zyra-api-gateway/internals/models"
	"github.com/AthulKrishna2501/zyra-api-gateway/pkg/config"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/webhook"
)

func PayMasterOfCeremony(ctx *gin.Context, c pb.ClientServiceClient) {
	var body models.MasterOfCeremonyRequest

	clientID, exists := ctx.Get("client_id")
	log.Print("Client ID in token:", clientID)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Client ID not found in token"})
		return
	}

	clientIDStr, ok := clientID.(string)

	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid client ID format"})
		return
	}

	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Fields cannot be empty"})
		return
	}

	if body.Method != "stripe" && body.Method != "razorpay" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Please enter a valid payment method"})
		return
	}

	parsedClientID, err := uuid.Parse(clientIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Client ID UUID format"})
		return
	}

	grpcReq := &pb.MasterOfCeremonyRequest{
		UserId: parsedClientID.String(),
		Method: body.Method,
	}

	res, err := c.GetMasterOfCeremony(ctx, grpcReq)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, &res.Url)
}

func HandleStripeWebhook(ctx *gin.Context, c pb.ClientServiceClient, cfg config.Config) {
	stripe.Key = cfg.STRIPE_SECRET_KEY
	endpointSecret := cfg.STRIPE_WEBHOOK_SECRET

	log.Print("Stripe API KEY :", stripe.Key)
	log.Print("Endpoint Secret :", endpointSecret)

	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}
	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	signatureHeader := ctx.GetHeader("Stripe-Signature")
	fmt.Println("🔹 Received Stripe Signature:", signatureHeader)

	event, err := webhook.ConstructEvent(body, signatureHeader, endpointSecret)
	if err != nil {
		fmt.Println("❌ Webhook verification failed:", err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Webhook signature verification failed"})
		return
	}

	fmt.Println("✅ Received event:", event.Type)

	_, err = c.HandleStripeEvent(context.Background(), &pb.StripeWebhookRequest{
		EventType: event.Type,
		Payload:   string(body),
	})
	if err != nil {
		fmt.Println("❌ Error forwarding event to Client Service:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process event"})
		return
	}

	ctx.Status(http.StatusOK)
}
