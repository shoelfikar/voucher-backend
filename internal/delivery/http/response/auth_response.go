package response

// LoginResponse represents the login response
type LoginResponse struct {
	Token string   `json:"token"`
	User  UserInfo `json:"user"`
}

// UserInfo represents user information in response
type UserInfo struct {
	Email string `json:"email"`
}
