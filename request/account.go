package request

//	@description	CreateAccountRequest request
//	@description	with user id and account type
type CreateAccountRequest struct {
	// User ID
	UserID string `json:"userID" binding:"required,uuid" example:"425129d3-72b3-4c64-8556-fe7da1889981"`
	// Account type. One of the following: 'checking', 'saving'
	Type string `json:"type" binding:"required" example:"checking" enums:"checking,saving"`
} //@name CreateAccountRequest

//	@description	MonetaryRequest request
//	@description	with user id and amount to deposit
type MonetaryRequest struct {
	// User ID
	UserID string `json:"userID" binding:"required,uuid" example:"425129d3-72b3-4c64-8556-fe7da1889981"`
	// Amount to deposit or withdraw
	Amount float64 `json:"amount" binding:"required" example:"45.12" minimum:"1" validate:"required"`
} //@name MonetaryRequest
