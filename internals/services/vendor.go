package services

import (
	"net/http"

	pb "github.com/AthulKrishna2501/proto-repo/vendor"
	"github.com/AthulKrishna2501/zyra-api-gateway/internals/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequestCategory(ctx *gin.Context, c pb.VendorSeviceClient) {
	var body models.RequestCategoryRequest

	vendorID, exists := ctx.Get("vendor_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Vendor ID not found in token"})
		return
	}

	vendorIDStr, ok := vendorID.(string)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid vendor ID format"})
		return
	}

	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Fields cannot be empty"})
		return
	}

	parsedVendorID, err := uuid.Parse(vendorIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vendor ID UUID format"})
		return
	}

	grpcReq := &pb.RequestCategoryRequest{
		VendorId:   parsedVendorID.String(),
		CategoryId: body.CategoryId,
	}

	res, err := c.RequestCategory(ctx, grpcReq)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(int(res.Status), &res)

}

func ListCategory(ctx *gin.Context, c pb.VendorSeviceClient) {
	grpcReq := &pb.ListCategoryRequest{}

	res, err := c.ListCategory(ctx, grpcReq)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, &res)

}
