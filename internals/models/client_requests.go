package models

type MasterOfCeremonyRequest struct {
	Method string `json:"method"`
}

type VerifyPaymentRequest struct {
	SessionID string `json:"session_id"`
}
