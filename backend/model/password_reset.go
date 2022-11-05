package model

type PasswordReset struct {
	Email   string `json:"email"`
	Token   string `json:"token"`
	Expired int64  `json:"expired"`
}

type UpdateResetPasswordUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

type CreatePasswordResetRequest struct {
	Email string `json:"email"`
}
