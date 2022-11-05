package comment

type CreateCommentRequest struct {
	UserOrderId int64  `json:"user_order_id"`
	ProductCode string `json:"product_code"`
	Description string `json:"description"`
	Tag         string `json:"tag"`
	Rating      int32  `json:"rating" binding:"required"`
}
