package request

type CreateAccount struct {
	UserID string `json:"userID"`
	Type   string `json:"type"`
}
