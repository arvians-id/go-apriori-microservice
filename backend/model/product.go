package model

import "time"

type Product struct {
	IdProduct   int64     `json:"id_product"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	Price       int32     `json:"price"`
	Category    string    `json:"category"`
	IsEmpty     bool      `json:"is_empty"`
	Mass        int32     `json:"mass"`
	Image       *string   `json:"image"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ProductRecommendation struct {
	AprioriId          int64   `json:"apriori_id"`
	AprioriCode        string  `json:"apriori_code"`
	AprioriItem        string  `json:"apriori_item"`
	AprioriDiscount    float64 `json:"apriori_discount"`
	AprioriDescription *string `json:"apriori_description"`
	AprioriImage       *string `json:"apriori_image"`
	ProductTotalPrice  int32   `json:"product_total_price"`
	PriceDiscount      int32   `json:"price_discount"`
	Mass               int32   `json:"mass"`
}

type GetProductNameTransactionResponse struct {
	ProductName []string `json:"product_name"`
}
