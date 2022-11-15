package model

import (
	"github.com/arvians-id/go-apriori-microservice/services/user-order-service/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type Apriori struct {
	IdApriori   int64      `json:"id_apriori"`
	Code        string     `json:"code"`
	Item        string     `json:"item"`
	Discount    float32    `json:"discount"`
	Support     float32    `json:"support"`
	Confidence  float32    `json:"confidence"`
	RangeDate   string     `json:"range_date"`
	IsActive    bool       `json:"is_active"`
	Description *string    `json:"description"`
	Mass        int32      `json:"mass"`
	Image       *string    `json:"image"`
	CreatedAt   time.Time  `json:"created_at"`
	UserOrder   *UserOrder `json:"user_order"`
}

func (apriori *Apriori) ToProtoBuff() *pb.Apriori {
	return &pb.Apriori{
		IdApriori:   apriori.IdApriori,
		Code:        apriori.Code,
		Item:        apriori.Item,
		Discount:    apriori.Discount,
		Support:     apriori.Support,
		Confidence:  apriori.Confidence,
		RangeDate:   apriori.RangeDate,
		IsActive:    apriori.IsActive,
		Description: apriori.Description,
		Mass:        apriori.Mass,
		Image:       apriori.Image,
		CreatedAt:   timestamppb.New(apriori.CreatedAt),
		//UserOrder:   apriori.UserOrder.ToProtoBuff(),
	}
}

type GenerateApriori struct {
	ItemSet     []string `json:"item_set"`
	Support     float32  `json:"support"`
	Iterate     int32    `json:"iterate"`
	Transaction int32    `json:"transaction"`
	Confidence  float32  `json:"confidence"`
	Discount    float32  `json:"discount"`
	Description string   `json:"description"`
	RangeDate   string   `json:"range_date"`
}
