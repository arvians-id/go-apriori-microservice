package model

import "time"

type Product struct {
	IdProduct   int       `json:"id_product"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	Price       int       `json:"price"`
	Category    string    `json:"category"`
	IsEmpty     bool      `json:"is_empty"`
	Mass        int       `json:"mass"`
	Image       *string   `json:"image"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ProductRecommendation struct {
	AprioriId          int     `json:"apriori_id"`
	AprioriCode        string  `json:"apriori_code"`
	AprioriItem        string  `json:"apriori_item"`
	AprioriDiscount    float64 `json:"apriori_discount"`
	AprioriDescription *string `json:"apriori_description"`
	AprioriImage       *string `json:"apriori_image"`
	ProductTotalPrice  int     `json:"product_total_price"`
	PriceDiscount      int     `json:"price_discount"`
	Mass               int     `json:"mass"`
}

type RatingFromComment struct {
	Rating        int `json:"rating"`
	ResultRating  int `json:"result_rating"`
	ResultComment int `json:"result_comment"`
}

type CreateProductRequest struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	Category    string `json:"category"`
	Mass        int    `json:"mass"`
	Image       string `json:"image"`
}

type UpdateProductRequest struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	Category    string `json:"category"`
	IsEmpty     bool   `json:"is_empty"`
	Mass        int    `json:"mass"`
	Image       string `json:"image"`
}

type GetProductNameTransactionResponse struct {
	ProductName []string `json:"product_name"`
}
