package request

import "time"

type CreatePromoRequest struct {
	Code          string    `json:"code" binding:"required"`
	DiscountType  string    `json:"discount_type" binding:"required,oneof=percentage fixed"`
	DiscountValue float64   `json:"discount_value" binding:"required"`
	ValidUntil    time.Time `json:"valid_until" binding:"required"`
}

type UpdatePromoRequest struct {
	Code          string    `json:"code"`
	DiscountType  string    `json:"discount_type" binding:"omitempty,oneof=percentage fixed"`
	DiscountValue float64   `json:"discount_value"`
	ValidUntil    time.Time `json:"valid_until"`
}
