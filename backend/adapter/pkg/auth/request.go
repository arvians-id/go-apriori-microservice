package auth

type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Address  string `json:"address"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type GetRefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type GetUserCredentialRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateResetPasswordUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

type CreatePasswordResetRequest struct {
	Email string `json:"email"`
}
