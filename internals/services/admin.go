package services

import (
	"log"
	"net/http"

	pb "github.com/AthulKrishna2501/proto-repo/admin"
	"github.com/AthulKrishna2501/zyra-api-gateway/internals/models"

	"github.com/gin-gonic/gin"
)

func ApproveRejectCategory(ctx *gin.Context, c pb.AdminServiceClient) {
	var body models.CategoryApproveReject

	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	log.Printf("API Gateway: Forwarding request to Admin Service - VendorID=%s, CategoryID=%s, Status=%s", body.VendorID, body.CategoryID, body.Status)

	grpcReq := &pb.ApproveRejectCategoryRequest{
		VendorId:   body.VendorID,
		CategoryId: body.CategoryID,
		Status:     body.Status,
	}

	res, err := c.ApproveRejectCategory(ctx, grpcReq)
	if err != nil {
		log.Printf("API Gateway: gRPC error when calling Admin Service: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func 	BlockUser(ctx *gin.Context, c pb.AdminServiceClient) {
	var body struct {
		UserID string `json:"user_id"`
	}

	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	grpcReq := &pb.BlockUnblockUserRequest{UserId: body.UserID}

	res, err := c.BlockUser(ctx, grpcReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func UnblockUser(ctx *gin.Context, c pb.AdminServiceClient) {
	var body struct {
		UserID string `json:"user_id"`
	}

	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	grpcReq := &pb.BlockUnblockUserRequest{UserId: body.UserID}

	res, err := c.UnblockUser(ctx, grpcReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func ListUsers(ctx *gin.Context, c pb.AdminServiceClient) {
	grpcReq := &pb.ListUsersRequest{}

	res, err := c.ListUsers(ctx, grpcReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, res)
}
