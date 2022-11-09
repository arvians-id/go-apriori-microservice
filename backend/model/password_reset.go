package model

import "github.com/arvians-id/go-apriori-microservice/adapter/pb"

type PasswordReset struct {
	Email   string `json:"email"`
	Token   string `json:"token"`
	Expired int64  `json:"expired"`
}

func (passwordReset *PasswordReset) ToProtoBuff() *pb.PasswordReset {
	return &pb.PasswordReset{
		Email:   passwordReset.Email,
		Token:   passwordReset.Token,
		Expired: passwordReset.Expired,
	}
}
