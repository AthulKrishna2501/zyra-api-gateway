package services

import (
	"net/http"

	pb "github.com/AthulKrishna2501/proto-repo/vendor"
	"github.com/AthulKrishna2501/zyra-api-gateway/internals/models"
	"github.com/AthulKrishna2501/zyra-api-gateway/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
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
		VendorId:     parsedVendorID.String(),
		CategoryName: body.CategoryName,
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

func VendorProfile(ctx *gin.Context, c pb.VendorSeviceClient) {
	vendorID, ok := utils.GetVendorID(ctx)
	if !ok {
		return
	}

	grpcReq := &pb.VendorProfileRequest{
		VendorId: vendorID.String(),
	}

	res, err := c.VendorProfile(ctx, grpcReq)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, &res)

}

func UpdateProfile(ctx *gin.Context, c pb.VendorSeviceClient) {
	var body models.UpdateVendorProfileRequest
	vendorID, ok := utils.GetVendorID(ctx)
	if !ok {
		return
	}

	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Fields cannot be empty"})
		return
	}

	grpcReq := &pb.UpdateVendorProfileRequest{
		VendorId:     vendorID.String(),
		FirstName:    &body.FirstName,
		LastAme:      &body.LastName,
		PhoneNumber:  &body.PhoneNumber,
		ProfileImage: &body.ProfileImage,
		Bio:          &body.Bio,
	}

	res, err := c.UpdateProfile(ctx, grpcReq)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, &res)

}

func CreateService(ctx *gin.Context, c pb.VendorSeviceClient) {
	var body models.CreateServiceRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := utils.ValidateServiceRequest(body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	vendorID, ok := utils.GetVendorID(ctx)
	if !ok {
		return
	}

	availableDate := timestamppb.New(body.AvailableDate)

	res, err := c.CreateService(ctx, &pb.CreateServiceRequest{
		VendorId:           vendorID.String(),
		YearOfExperience:   body.YearOfExperience,
		AvailableDates:     []*timestamppb.Timestamp{availableDate},
		ServiceDescription: body.ServiceDescription,
		ServiceDuration:    body.ServiceDuration,
		ServicePrice:       body.ServicePrice,
		ServiceTitle:       body.ServiceTitle,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, &res)
}

func UpdateService(ctx *gin.Context, c pb.VendorSeviceClient) {
	var req models.UpdateServiceRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := utils.ValidateUpdateRequest(req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	availableDate := timestamppb.New(req.AvailableDate)

	grpcReq := &pb.UpdateServiceRequest{
		ServiceId:           req.ServiceID.String(),
		ServiceTitle:        req.ServiceTitle,
		YearOfExperience:    req.YearOfExperience,
		ServiceDescription:  req.ServiceDescription,
		AvailableDates:      []*timestamppb.Timestamp{availableDate},
		CancellationPolicy:  req.CancellationPolicy,
		TermsAndConditions:  req.TermsAndConditions,
		ServiceDuration:     req.ServiceDuration,
		ServicePrice:        req.ServicePrice,
		AdditionalHourPrice: &req.AdditionalHourPrice,
	}
	res, err := c.UpdateService(ctx, grpcReq)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, &res)
}

func ChangePassword(ctx *gin.Context, c pb.VendorSeviceClient) {
	var req models.ChangePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Fields cannot be empty"})
		return
	}

	vendorID, ok := utils.GetVendorID(ctx)
	if !ok {
		return
	}

	if len(req.NewPassword) < 8 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Password should be atleast 8 characters"})
		return
	}

	if req.NewPassword != req.ConfirmPassword {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "New password and confirm password do not match"})
		return
	}

	grpcReq := &pb.ChangePasswordRequest{
		VendorId:        vendorID.String(),
		CurrentPassword: req.CurrentPassword,
		NewPassword:     req.NewPassword,
		ConfirmPassword: req.ConfirmPassword,
	}

	res, err := c.ChangePassword(ctx, grpcReq)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": res.Message})
}

func VendorDashBoard(ctx *gin.Context, c pb.VendorSeviceClient) {
	vendorID, ok := utils.GetVendorID(ctx)
	if !ok {
		return
	}

	grpReq := &pb.GetVendorDashboardRequest{
		VendorId: vendorID.String(),
	}

	res, err := c.GetVendorDashboard(ctx, grpReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return

	}

	ctx.JSON(http.StatusOK, gin.H{"message": res})
}

func GetServices(ctx *gin.Context, c pb.VendorSeviceClient) {
	vendorID, ok := utils.GetVendorID(ctx)
	if !ok {
		return
	}

	grpReq := &pb.GetVendorServicesRequest{
		VendorId: vendorID.String(),
	}

	res, err := c.GetVendorServices(ctx, grpReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return

	}

	ctx.JSON(http.StatusOK, gin.H{"message": res})

}
