package services

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	pb "github.com/AthulKrishna2501/proto-repo/client"
	"github.com/AthulKrishna2501/zyra-api-gateway/internals/constants"
	"github.com/AthulKrishna2501/zyra-api-gateway/internals/models"
	"github.com/AthulKrishna2501/zyra-api-gateway/pkg/config"
	"github.com/AthulKrishna2501/zyra-api-gateway/pkg/validator"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/webhook"
	"google.golang.org/protobuf/types/known/timestamppb"
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

func HostEvent(ctx *gin.Context, c pb.ClientServiceClient) {
	var req models.CreateEventRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request"})
		return
	}

	clientID, exists := ctx.Get("user_ids")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Client ID not found in token"})
		return
	}

	clientIDStr, ok := clientID.(string)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid client ID format"})
		return
	}

	eventDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event date"})
		return
	}

	startTime, err := time.Parse("15:04", req.EventDetails.StartTime)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start time"})
		return
	}

	endTime, err := time.Parse("15:04", req.EventDetails.EndTime)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end time"})
		return
	}

	eventID := uuid.New()
	grpcReq := &pb.CreateEventRequest{
		EventId:  eventID.String(),
		Title:    req.Title,
		Date:     timestamppb.New(eventDate),
		HostedBy: clientIDStr,
		Location: &pb.Location{
			Address:   req.Location.Address,
			City:      req.Location.City,
			Country:   req.Location.Country,
			Latitude:  req.Location.Lat,
			Longitude: req.Location.Lng,
		},

		EventDetails: &pb.EventDetails{
			EventId:        eventID.String(),
			Description:    req.EventDetails.Description,
			StartTime:      timestamppb.New(startTime),
			EndTime:        timestamppb.New(endTime),
			PosterImage:    req.EventDetails.PosterImage,
			PricePerTicket: int32(req.EventDetails.PricePerTicket),
			TicketLimit:    int32(req.EventDetails.TicketLimit),
		},
	}

	res, err := c.CreateEvent(ctx, grpcReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, res)

}

func EditEvent(ctx *gin.Context, c pb.ClientServiceClient) {
	var req models.EditEventRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request"})
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

	eventDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event date"})
		return
	}

	startTime, err := time.Parse("15:04", req.EventDetails.StartTime)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start time"})
		return
	}

	endTime, err := time.Parse("15:04", req.EventDetails.EndTime)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end time"})
		return
	}

	grpcReq := &pb.EditEventRequest{
		EventId:  req.EventId,
		Title:    req.Title,
		Date:     timestamppb.New(eventDate),
		HostedBy: clientIDStr,
		Location: &pb.Location{
			Address:   req.Location.Address,
			City:      req.Location.City,
			Country:   req.Location.Country,
			Latitude:  req.Location.Lat,
			Longitude: req.Location.Lng,
		},

		EventDetails: &pb.EventDetails{
			EventId:        req.EventId,
			Description:    req.EventDetails.Description,
			StartTime:      timestamppb.New(startTime),
			EndTime:        timestamppb.New(endTime),
			PosterImage:    req.EventDetails.PosterImage,
			PricePerTicket: int32(req.EventDetails.PricePerTicket),
			TicketLimit:    int32(req.EventDetails.TicketLimit),
		},
	}

	res, err := c.EditEvent(ctx, grpcReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func GetClientProfile(ctx *gin.Context, c pb.ClientServiceClient) {
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

	grpcReq := &pb.GetClientProfileRequest{
		ClientId: clientIDStr,
	}

	res, err := c.GetClientProfile(ctx, grpcReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to fetch client profile", "details": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"clientId":     res.ClientId,
			"firstName":    res.FirstName,
			"lastName":     res.LastName,
			"email":        res.Email,
			"profileImage": res.ProfileImage,
			"phoneNumber":  res.PhoneNumber,
		},
	})
}

func EditClientProfile(ctx *gin.Context, c pb.ClientServiceClient) {
	var req models.EditClientProfileRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
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

	if !validator.ValidatePhone(req.PhoneNumber) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "phone number should be 10 digits"})
		return
	}

	grpcReq := &pb.EditClientProfileRequest{
		ClientId:     clientIDStr,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		ProfileImage: req.ProfileImage,
		PhoneNumber:  req.PhoneNumber,
	}

	res, err := c.EditClientProfile(ctx, grpcReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to edit client profile", "details": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": res.Message,
	})
}

func ResetPassword(ctx *gin.Context, c pb.ClientServiceClient) {
	var req models.ResetPasswordRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if len(req.NewPassword) < constants.PasswordMinLength {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "password must be at least 8 characters long"})
		return
	}

	if req.NewPassword != req.ConfirmPassword {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Password and confirm password do not match"})
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

	grpcReq := &pb.ResetPasswordRequest{
		ClientId:        clientIDStr,
		CurrentPassword: req.CurrentPassword,
		NewPassword:     req.NewPassword,
		ConfirmPassword: req.ConfirmPassword,
	}

	res, err := c.ResetPassword(ctx, grpcReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to reset password", "details": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": res.Message,
	})
}
