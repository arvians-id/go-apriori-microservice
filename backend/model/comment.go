package model

import "time"

type Comment struct {
	IdComment   int64     `json:"id_comment"`
	UserOrderId int64     `json:"user_order_id"`
	ProductCode string    `json:"product_code"`
	Description *string   `json:"description"`
	Tag         *string   `json:"tag"`
	Rating      int32     `json:"rating"`
	CreatedAt   time.Time `json:"created_at"`
}
