package request

type CreateAccount struct {
	UserID string `json:"userID" binding:"required,uuid"`
	Type   string `json:"type" binding:"required"`
}
