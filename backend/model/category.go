package model

import (
	"github.com/arvians-id/go-apriori-microservice/adapter/pkg/category/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type Category struct {
	IdCategory int64     `json:"id_category"`
	Name       string    `json:"name"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (category *Category) ToProtoBuff() *pb.Category {
	return &pb.Category{
		IdCategory: category.IdCategory,
		Name:       category.Name,
		CreatedAt:  timestamppb.New(category.CreatedAt),
		UpdatedAt:  timestamppb.New(category.UpdatedAt),
	}
}
