package model

import "time"

type UpdateUserRequest struct {
	IdUser   *int   `json:"id_user"`
	Role     int    `json:"role"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Address  string `json:"address"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type User struct {
	IdUser       int             `json:"id_user"`
	Role         int             `json:"role"`
	Name         string          `json:"name"`
	Email        string          `json:"email"`
	Address      string          `json:"address"`
	Phone        string          `json:"phone"`
	Password     string          `json:"password"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
	Notification []*Notification `json:"notification"`
	Payment      []*Payment      `json:"payment"`
}

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

type TokenJwt struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
