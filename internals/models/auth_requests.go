package models

type RegisterRequestBody struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"required"`
}

type LoginRequestBody struct {
	Email    string `json:"email" binding:"required"`
	Role     string `json:"role" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type OTPRequestBody struct {
	Email string `json:"email" binding:"required"`
	Role  string `json:"role" binding:"required"`
}

type VerifyOTPBody struct {
	Email string `json:"email" binding:"required"`
	OTP   string `json:"otp" binding:"required"`
}
