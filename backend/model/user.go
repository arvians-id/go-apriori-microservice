package model

import "time"

type User struct {
	IdUser    int64     `json:"id_user"`
	Role      int32     `json:"role"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Address   string    `json:"address"`
	Phone     string    `json:"phone"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TokenJwt struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
