package models

import (
	"github.com/AthulKrishna2501/zyra-client-service/internals/core/models"
)

type GenericBookingRequest struct {
	Method      string            `json:"method"`
	ServiceType string            `json:"service_type"`
	Metadata    map[string]string `json:"metadata"`
}

type VerifyPaymentRequest struct {
	SessionID string `json:"session_id"`
}

type CreateEventRequest struct {
	Title        string          `json:"title"`
	Date         string          `json:"date"`
	Location     models.Location `json:"location"`
	EventDetails struct {
		Description    string `json:"description"`
		StartTime      string `json:"start_time"`
		EndTime        string `json:"end_time"`
		PosterImage    string `json:"poster_image"`
		PricePerTicket int    `json:"price_per_ticket"`
		TicketLimit    int    `json:"ticket_limit"`
	} `json:"event_details"`
}

type EditEventRequest struct {
	EventId      string          `json:"event_id"`
	Title        string          `json:"title"`
	Date         string          `json:"date"`
	Location     models.Location `json:"location"`
	EventDetails struct {
		Description    string `json:"description"`
		StartTime      string `json:"start_time"`
		EndTime        string `json:"end_time"`
		PosterImage    string `json:"poster_image"`
		PricePerTicket int    `json:"price_per_ticket"`
		TicketLimit    int    `json:"ticket_limit"`
	} `json:"event_details"`
}

type EditClientProfileRequest struct {
	FirstName    string `json:"first_name" binding:"required"`
	LastName     string `json:"last_name" binding:"required"`
	Place        string `json:"place" binding:"required"`
	ProfileImage string `json:"profile_image" binding:"required"`
	PhoneNumber  string `json:"phone_number" binding:"required"`
}

type ResetPasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
}

type BookVendorRequest struct {
	VendorId  string `json:"vendor_id" binding:"required"`
	ServiceId string `json:"service_id" binding:"required"`
	Date      string `json:"date" binding:"required"`
}

type ReviewRatingsRequest struct {
	VendorID string  `json:"vendor_id" binding:"required"`
	Rating   float64 `json:"rating" binding:"required"`
	Review   string  `json:"review"`
}

type EditReviewRatingsRequest struct {
	ReviewID string  `json:"review_id" binding:"required"`
	Rating   float64 `json:"rating" binding:"required"`
	Review   string  `json:"review"`
}

type DeleteReviewRequest struct {
	ReviewID string `json:"review_id"`
}

type CompleteVendorBookingRequest struct {
	BookingID string `json:"booking_id"`
	Status    string `json:"status"`
}

type CancelVendorBookingRequest struct {
	BookingID string `json:"booking_id"`
}

type CancelEventRequest struct {
	EventID string `json:"event_id"`
}

type FundReleaseRequest struct {
	EventID string `json:"event_id"`
}
