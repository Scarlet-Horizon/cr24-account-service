package request

//	@description	User request
//	@description	with user id
type User struct {
	// User ID
	UserID string `json:"userID" binding:"required,uuid" example:"425129d3-72b3-4c64-8556-fe7da1889981"`
} //	@name	User

//	@description	AccountRequest request
//	@description	with user id and account type
type AccountRequest struct {
	User
	// Account type. One of the following: 'checking', 'saving'
	Type string `json:"type" binding:"required" example:"checking" enums:"checking,saving"`
} //@name AccountRequest

//	@description	MonetaryRequest request
//	@description	with user id and amount to deposit
type MonetaryRequest struct {
	User
	// Amount to deposit or withdraw
	Amount float64 `json:"amount" binding:"required" example:"45.12" minimum:"1" validate:"required"`
} //	@name	MonetaryRequest

//	@description	CloseRequest request
//	@description	with user id
type CloseRequest struct {
	User
} //	@name	CloseRequest
