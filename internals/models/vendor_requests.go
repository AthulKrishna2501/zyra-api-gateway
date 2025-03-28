package models

import (
	"time"

	"github.com/google/uuid"
)

type RequestCategoryRequest struct {
	CategoryId string `json:"category_id"`
}

type UpdateVendorProfileRequest struct {
	FirstName    string `json:"first_name,omitempty"`
	LastName     string `json:"last_name,omitempty"`
	Place        string `json:"place,omitempty"`
	ProfileImage string `json:"profile_image,omitempty"`
	PhoneNumber  string `json:"phone_number,omitempty"`
	Bio          string `json:"bio,omitempty"`
}

type CreateServiceRequest struct {
	ServiceTitle        string    `json:"service_title" binding:"required"`
	YearOfExperience    int32     `json:"year_of_experience" binding:"required"`
	ServiceDescription  string    `json:"service_description" binding:"required"`
	AvailableDate       time.Time `json:"available_date" binding:"required"`
	CancellationPolicy  string    `json:"cancellation_policy,omitempty"`
	TermsAndConditions  string    `json:"terms_and_conditions,omitempty"`
	ServiceDuration     int32     `json:"service_duration" binding:"required"`
	ServicePrice        int32     `json:"service_price" binding:"required"`
	AdditionalHourPrice int32     `json:"additional_hour_price,omitempty"`
}
type UpdateServiceRequest struct {
	ServiceID           uuid.UUID `json:"service_id" binding:"required"`
	ServiceTitle        string    `json:"service_title" binding:"required"`
	YearOfExperience    int32     `json:"year_of_experience" binding:"required"`
	ServiceDescription  string    `json:"service_description" binding:"required"`
	AvailableDate       time.Time `json:"available_date" binding:"required"`
	CancellationPolicy  string    `json:"cancellation_policy,omitempty"`
	TermsAndConditions  string    `json:"terms_and_conditions,omitempty"`
	ServiceDuration     int32     `json:"service_duration" binding:"required"`
	ServicePrice        int32     `json:"service_price" binding:"required"`
	AdditionalHourPrice int32     `json:"additional_hour_price,omitempty"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
}
