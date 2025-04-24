package utils

import (
	"errors"
	"net/http"

	"github.com/AthulKrishna2501/zyra-api-gateway/internals/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetVendorID(ctx *gin.Context) (uuid.UUID, bool) {
	vendorID, exists := ctx.Get("vendor_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Vendor ID not found in token"})
		ctx.Abort()
		return uuid.Nil, false
	}

	vendorIDStr, ok := vendorID.(string)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid vendor ID format"})
		ctx.Abort()
		return uuid.Nil, false
	}

	parsedVendorID, err := uuid.Parse(vendorIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vendor ID UUID format"})
		ctx.Abort()
		return uuid.Nil, false
	}

	return parsedVendorID, true
}

func ValidateServiceRequest(req models.CreateServiceRequest) error {
	if len(req.ServiceTitle) == 0 {
		return errors.New("service_title cannot be empty")
	}

	if req.YearOfExperience <= 0 {
		return errors.New("year_of_experience must be greater than 0")
	}

	if len(req.ServiceDescription) == 0 {
		return errors.New("service_description cannot be empty")
	}

	if req.ServiceDuration <= 0 {
		return errors.New("service_duration must be greater than 0")
	}

	if req.ServicePrice <= 0 {
		return errors.New("service_price must be greater than 0")
	}

	if req.AdditionalHourPrice < 0 {
		return errors.New("additional_hour_price must be greater than or equal to 0")
	}

	return nil
}

func ValidateUpdateRequest(req models.UpdateServiceRequest) error {
	if len(req.ServiceTitle) == 0 {
		return errors.New("service_title cannot be empty")
	}

	if req.YearOfExperience <= 0 {
		return errors.New("year_of_experience must be greater than 0")
	}

	if len(req.ServiceDescription) == 0 {
		return errors.New("service_description cannot be empty")
	}

	if req.ServiceDuration <= 0 {
		return errors.New("service_duration must be greater than 0")
	}

	if req.ServicePrice <= 0 {
		return errors.New("service_price must be greater than 0")
	}

	if req.AdditionalHourPrice < 0 {
		return errors.New("additional_hour_price must be greater than or equal to 0")
	}

	return nil
}

func ValidateReviewRating(ratings float64, review string) error {
	if ratings > 10 || ratings < 0 {
		return errors.New("please enter a valid rating")
	}

	if review == "" {
		return errors.New("please provide a review before submitting")
	}

	return nil
}
