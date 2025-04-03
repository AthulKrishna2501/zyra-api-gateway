package services

import (
	"bytes"
	"context"
	"errors"
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

func PayMasterOfCeremony(ctx *gin.Context, c pb.ClientServiceClient) error {
	var body models.MasterOfCeremonyRequest

	clientID, exists := ctx.Get("client_id")
	log.Print("Client ID in token:", clientID)
	if !exists {
		return errors.New("client ID not found in token")
	}

	clientIDStr, ok := clientID.(string)
	if !ok {
		return errors.New("invalid client ID format")
	}

	if err := ctx.BindJSON(&body); err != nil {
		return errors.New("fields cannot be empty")
	}

	if body.Method != "stripe" && body.Method != "razorpay" {
		return errors.New("please enter a valid payment method")
	}

	parsedClientID, err := uuid.Parse(clientIDStr)
	if err != nil {
		return errors.New("invalid Client ID UUID format")
	}

	grpcReq := &pb.MasterOfCeremonyRequest{
		UserId: parsedClientID.String(),
		Method: body.Method,
	}

	res, err := c.GetMasterOfCeremony(ctx, grpcReq)
	if err != nil {
		return fmt.Errorf("gRPC request failed: %v", err)
	}

	ctx.JSON(http.StatusOK, gin.H{"url": res.Url})
	return nil
}

func HandleStripeWebhook(ctx *gin.Context, c pb.ClientServiceClient, cfg config.Config) {
	stripe.Key = cfg.STRIPE_SECRET_KEY
	endpointSecret := cfg.STRIPE_WEBHOOK_SECRET

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
		fmt.Println("Webhook verification failed:", err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Webhook signature verification failed"})
		return
	}

	fmt.Println("Received event:", event.Type)

	_, err = c.HandleStripeEvent(context.Background(), &pb.StripeWebhookRequest{
		EventType: event.Type,
		Payload:   string(body),
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process event"})
		return
	}

	ctx.Status(http.StatusOK)
}

func VerifyPayment(ctx *gin.Context, c pb.ClientServiceClient) {
	sessionID := ctx.Query("session_id")
	log.Print("session ID in api gateway :", sessionID)
	if sessionID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Session ID is required"})
		return
	}

	clientID, exists := ctx.Get("client_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Client ID not found in token"})
		return
	}

	clientIDStr, ok := clientID.(string)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid client ID format"})
		return
	}

	grpcReq := &pb.VerifyPaymentRequest{
		UserId:    clientIDStr,
		SessionId: sessionID,
	}

	log.Printf("Sending gRPC request for verifyPayment with session_id %s :", sessionID)

	res, err := c.VerifyPayment(ctx, grpcReq)

	if err != nil {
		log.Print("Error in verifying payment:", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Payment verification failed", "details": err.Error()})
		return
	}

	log.Printf("gRPC response: %+v", res)

	if res.Success {
		ctx.JSON(http.StatusOK, gin.H{
			"message": res.Message,
		})
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": res.Message,
		})
	}
}
