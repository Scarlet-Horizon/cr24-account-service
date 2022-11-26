package response

type ErrorResponse struct {
	// Error description.
	Error string `json:"error" example:"invalid user id"`
} //@name ErrorResponse
