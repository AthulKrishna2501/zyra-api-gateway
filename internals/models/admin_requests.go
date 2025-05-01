package models

type CategoryApproveReject struct {
	VendorID   string `json:"vendor_id"`
	CategoryID string `json:"category_id"`
	Status     string `json:"status"`
}

type AddCategoryRequest struct {
	CategoryName string `json:"category_name" binding:"required"`
}

type ApproveFundReleaseRequest struct {
	RequestID string `json:"request_id"`
	Status    string `json:"status"`
}
