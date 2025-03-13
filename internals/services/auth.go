package services

import (
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	pb "github.com/AthulKrishna2501/proto-repo/auth"
	"github.com/AthulKrishna2501/zyra-api-gateway/internals/events"
	"github.com/AthulKrishna2501/zyra-api-gateway/internals/models"
	"github.com/AthulKrishna2501/zyra-api-gateway/pkg/validator"
	"github.com/gin-gonic/gin"
)

func Register(ctx *gin.Context, c pb.AuthServiceClient, mq *events.RabbitMq) {
	body := models.RegisterRequestBody{}

	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Fields cannot be empty"})
		return
	}

	if err := validator.ValidateSignup(body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	grpcReq := &pb.RegisterRequest{
		Name:     body.Name,
		Email:    body.Email,
		Password: body.Password,
		Role:     body.Role,
	}

	res, err := c.Register(ctx, grpcReq)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
		return
	}

	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)

	otp := rng.Intn(900000) + 100000
	otpStr := strconv.Itoa(otp)

	err = mq.PublishOTP(body.Email, otpStr)
	if err != nil {
		log.Println("Failed to Publish OTP ", err)
	} else {
		log.Printf("OTP %s published for email %s", otpStr, body.Email)
	}

	ctx.JSON(int(res.Status), &res)

}

func Login(ctx *gin.Context, c pb.AuthServiceClient) {
	body := models.LoginRequestBody{}

	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Fields cannot be empty"})
		return
	}

	if err := validator.ValidateLogin(body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	grpcReq := pb.LoginRequest{
		Email:    body.Email,
		Role:     body.Role,
		Password: body.Password,
	}

	res, err := c.Login(ctx, &grpcReq)

	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(int(res.Status), &res)

}

func SendOTP(ctx *gin.Context, c pb.AuthServiceClient) {
	body := models.OTPRequestBody{}

	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Fields cannot be empty"})
		return
	}

	err := validator.ValidateOTP(body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	grpcReq := pb.OTPRequest{
		Email: body.Email,
		Role:  body.Role,
	}

	res, err := c.SendOTP(ctx, &grpcReq)

	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(int(res.Status), &res)

}

func VerifyOTP(ctx *gin.Context, c pb.AuthServiceClient) {
	body := models.VerifyOTPBody{}
	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Fields cannot be empty"})
		return
	}

	err := validator.ValidateVerifyOTP(body)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	grpcReq := pb.VerifyOTPRequest{
		Email: body.Email,
		Otp:   body.OTP,
	}

	res, err := c.Verify(ctx, &grpcReq)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(int(res.Status), &res)

}

func Logout(ctx *gin.Context, c pb.AuthServiceClient) {
	grpcReq := pb.LogoutRequest{}
	res, err := c.Logout(ctx, &grpcReq)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(int(res.Status), &res)

}
