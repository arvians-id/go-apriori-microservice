package model

import "time"

type Comment struct {
	IdComment   int        `json:"id_comment"`
	UserOrderId int        `json:"user_order_id"`
	ProductCode string     `json:"product_code"`
	Description *string    `json:"description"`
	Tag         *string    `json:"tag"`
	Rating      int        `json:"rating"`
	CreatedAt   time.Time  `json:"created_at"`
	UserOrder   *UserOrder `json:"user_order"`
}

type CreateCommentRequest struct {
	UserOrderId int    `json:"user_order_id"`
	ProductCode string `json:"product_code"`
	Description string `json:"description"`
	Tag         string `json:"tag"`
	Rating      int    `json:"rating"`
}
