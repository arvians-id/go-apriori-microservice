package auth

type CreateUserRequest struct {
	Name     string `json:"name" binding:"required,max=20"`
	Email    string `json:"email" binding:"required,email,max=100"`
	Address  string `json:"address" binding:"required,max=100"`
	Phone    string `json:"phone" binding:"required,max=20"`
	Password string `json:"password" binding:"required,min=6"`
}

type GetRefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type GetUserCredentialRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UpdateResetPasswordUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Token    string `json:"token"`
}

type CreatePasswordResetRequest struct {
	Email string `json:"email" binding:"required,email"`
}
