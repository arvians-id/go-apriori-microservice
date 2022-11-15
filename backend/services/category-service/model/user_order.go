package model

import "github.com/arvians-id/go-apriori-microservice/services/category-service/pb"

type UserOrder struct {
	IdOrder        int64    `json:"id_order"`
	PayloadId      int64    `json:"payload_id"`
	Code           *string  `json:"code"`
	Name           *string  `json:"name"`
	Price          *int64   `json:"price"`
	Image          *string  `json:"image"`
	Quantity       *int32   `json:"quantity"`
	TotalPriceItem *int64   `json:"total_price_item"`
	Payment        *Payment `json:"payment"`
}

func (userOrder *UserOrder) ToProtoBuff() *pb.UserOrder {
	return &pb.UserOrder{
		IdOrder:        userOrder.IdOrder,
		PayloadId:      userOrder.PayloadId,
		Code:           userOrder.Code,
		Name:           userOrder.Name,
		Price:          userOrder.Price,
		Image:          userOrder.Image,
		Quantity:       userOrder.Quantity,
		TotalPriceItem: userOrder.TotalPriceItem,
		//Payment:        userOrder.Payment.ToProtoBuff(),
	}
}

func (userOrder *UserOrder) ToListProtoBuff() []*pb.UserOrder {
	return []*pb.UserOrder{
		{
			IdOrder:        userOrder.IdOrder,
			PayloadId:      userOrder.PayloadId,
			Code:           userOrder.Code,
			Name:           userOrder.Name,
			Price:          userOrder.Price,
			Image:          userOrder.Image,
			Quantity:       userOrder.Quantity,
			TotalPriceItem: userOrder.TotalPriceItem,
			//Payment:        userOrder.Payment.ToProtoBuff(),
		},
	}
}
