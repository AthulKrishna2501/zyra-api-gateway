package services

import (
	"net/http"
	"strings"

	pb "github.com/AthulKrishna2501/proto-repo/auth"
	"github.com/AthulKrishna2501/zyra-api-gateway/internals/models"
	"github.com/AthulKrishna2501/zyra-api-gateway/pkg/validator"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Register(ctx *gin.Context, c pb.AuthServiceClient) {
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
		ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(int(res.Status), &res)

}

func SendOTP(ctx *gin.Context, c pb.AuthServiceClient) {
	body := models.OTPRequestBody{}
	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Fields annot be empty"})
		return
	}

	if err := validator.ValidateOTP(body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	grpcReq := &pb.OTPRequest{
		Email: body.Email,
		Role:  body.Role,
	}

	res, err := c.SendOTP(ctx, grpcReq)
	if err != nil {
		grpcErrCode := status.Code(err)
		var httpStatus int
		switch grpcErrCode {
		case codes.Internal:
			httpStatus = http.StatusInternalServerError
		case codes.ResourceExhausted:
			httpStatus = http.StatusTooManyRequests
		default:
			httpStatus = http.StatusBadRequest

		}

		ctx.JSON(httpStatus, gin.H{"error": err.Error()})
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
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(int(res.Status), &res)

}

func ResendOTP(ctx *gin.Context, c pb.AuthServiceClient) {
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

	grpcReq := pb.ResendOTPRequest{
		Email: body.Email,
		Role:  body.Role,
	}

	res, err := c.ResendOTP(ctx, &grpcReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
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
		ctx.JSON(http.StatusForbidden, err.Error())
		return
	}

	ctx.JSON(int(res.Status), &res)

}

func RefreshToken(ctx *gin.Context, c pb.AuthServiceClient) {
	var body models.TokenRequest

	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Fields cannot be empty"})
		return
	}

	grpcReq := pb.RefreshTokenRequest{
		RefreshToken: body.RefreshToken,
	}

	res, err := c.RefreshToken(ctx, &grpcReq)

	if err != nil {
		ctx.JSON(http.StatusForbidden, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, &res)
}

func Logout(ctx *gin.Context, c pb.AuthServiceClient) {
	tokenString := ctx.GetHeader("Authorization")

	if tokenString == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
		return
	}

	tokenParts := strings.Split(tokenString, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
		return
	}
	token := tokenParts[1]

	grpcReq := &pb.LogoutRequest{
		AccessToken: token,
	}
	res, err := c.Logout(ctx, grpcReq)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(int(res.Status), &res)

}
