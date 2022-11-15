package model

import (
	"github.com/arvians-id/go-apriori-microservice/services/apriori-service/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type User struct {
	IdUser       int64           `json:"id_user"`
	Role         int32           `json:"role"`
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

func (user *User) ToProtoBuff() *pb.User {
	return &pb.User{
		IdUser:    user.IdUser,
		Role:      user.Role,
		Name:      user.Name,
		Email:     user.Email,
		Address:   user.Address,
		Phone:     user.Phone,
		Password:  user.Password,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}
}

type TokenJwt struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
