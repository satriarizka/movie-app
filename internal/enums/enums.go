package enums

// --- Role Enums ---
type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

// --- Transaction Status Enums ---
type TransactionStatus string

const (
	TransactionPending TransactionStatus = "pending"
	TransactionPaid    TransactionStatus = "paid"
	TransactionCancel  TransactionStatus = "cancelled"
	TransactionFailed  TransactionStatus = "failed"
)

// --- Payment Methods ---
const (
	PaymentCreditCard = "credit_card"
	PaymentEWallet    = "e_wallet"
	PaymentQRIS       = "qris"
)

// === Discount Types ===
const (
	DiscountTypePercentage = "percentage" // Misal: 10%
	DiscountTypeFixed      = "fixed"      // Misal: Potongan Rp 10.000
)
