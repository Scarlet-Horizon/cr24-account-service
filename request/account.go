package request

// AccountRequest godoc
// @Description AccountRequest with account type
type AccountRequest struct {
	// Account type. One of the following: 'checking', 'saving'
	Type string `json:"type" binding:"required" example:"checking" enums:"checking,saving"`
} //@Name AccountRequest

// MonetaryRequest godoc
// @Description	MonetaryRequest with amount to deposit
type MonetaryRequest struct {
	// Amount to deposit or withdraw
	Amount float64 `json:"amount" binding:"required" example:"45.12" minimum:"1" validate:"required"`
} //@Name MonetaryRequest
