package request

type PayTransactionRequest struct {
	PaymentMethod string `json:"payment_method" validate:"required,oneof=credit_card e_wallet qris"`
}
