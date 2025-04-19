package models

import "github.com/AthulKrishna2501/zyra-client-service/internals/core/models"

type MasterOfCeremonyRequest struct {
	Method string `json:"method"`
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
	VendorId string `json:"vendorId" binding:"required"`
	Service  string `json:"service" binding:"required"`
	Date     string `json:"date" binding:"required"`
}
