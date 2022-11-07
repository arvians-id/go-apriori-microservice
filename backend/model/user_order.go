package model

type UserOrder struct {
	IdOrder        int64   `json:"id_order"`
	PayloadId      int64   `json:"payload_id"`
	Code           *string `json:"code"`
	Name           *string `json:"name"`
	Price          *int64  `json:"price"`
	Image          *string `json:"image"`
	Quantity       *int32  `json:"quantity"`
	TotalPriceItem *int64  `json:"total_price_item"`
}
