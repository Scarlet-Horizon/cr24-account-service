package request

// CreateAccount request
type CreateAccount struct {
	// User ID
	UserID string `json:"userID" binding:"required,uuid" example:"425129d3-72b3-4c64-8556-fe7da1889981"`
	// Account type. One of the following: 'checking', 'saving'
	Type string `json:"type" binding:"required" example:"checking"`
}
