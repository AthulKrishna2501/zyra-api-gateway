package validator

import (
	"errors"

	"github.com/AthulKrishna2501/zyra-api-gateway/internals/constants"
	"github.com/AthulKrishna2501/zyra-api-gateway/internals/models"
)

func ValidateSignup(req models.RegisterRequestBody) error {
	if req.Name == "" {
		return errors.New("name is required")
	}

	if !constants.EmailRegex.MatchString(req.Email) {
		return errors.New("invalid email format")
	}

	if len(req.Password) < constants.PasswordMinLength {
		return errors.New("password must be at least 8 characters long")
	}

	return nil
}

func ValidateLogin(req models.LoginRequestBody) error {
	if req.Email == "" {
		return errors.New("email is required")
	}

	if req.Password == "" {
		return errors.New("password is required")
	}

	return nil
}

func ValidateOTP(req models.OTPRequestBody) error {
	if !constants.EmailRegex.MatchString(req.Email) {
		return errors.New("invalid email format")
	}

	if req.Role != "client" && req.Role != "vendor" {
		return errors.New("invalid role")
	}

	return nil
}

func ValidateVerifyOTP(req models.VerifyOTPBody) error {
	if !constants.EmailRegex.MatchString(req.Email) {
		return errors.New("invalid email format")
	}

	if len(req.OTP) < 6 {
		return errors.New("otp should be 6 digits")
	}

	return nil
}
