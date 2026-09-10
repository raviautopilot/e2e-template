package myservice_test

// LoginRequest represents the authentication request body.
type LoginRequest struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

// LoginResponse represents the authentication response payload.
type LoginResponse struct {
	Token       string                 `json:"token,omitempty"`
	AccessToken string                 `json:"access_token,omitempty"`
	Message     string                 `json:"message,omitempty"`
	User        map[string]interface{} `json:"user,omitempty"`
	Status      string                 `json:"status,omitempty"`
}
