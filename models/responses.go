package models

// Swagger types to show clear input / output messages.

// This type is only used for swagger documentation to show the API responses.
type ErrorResponse struct {
	Error string `json:"error"`
}

// Used for login endpoint
type TokenResponse struct {
	Token string `json:"token"`
}
