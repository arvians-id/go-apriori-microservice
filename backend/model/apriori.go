package model

import "time"

type Apriori struct {
	IdApriori   int        `json:"id_apriori"`
	Code        string     `json:"code"`
	Item        string     `json:"item"`
	Discount    float64    `json:"discount"`
	Support     float64    `json:"support"`
	Confidence  float64    `json:"confidence"`
	RangeDate   string     `json:"range_date"`
	IsActive    bool       `json:"is_active"`
	Description *string    `json:"description"`
	Mass        int        `json:"mass"`
	Image       *string    `json:"image"`
	CreatedAt   time.Time  `json:"created_at"`
	UserOrder   *UserOrder `json:"user_order"`
}

type CreateAprioriRequest struct {
	Item       string  `json:"item"`
	Discount   float64 `json:"discount"`
	Support    float64 `json:"support"`
	Confidence float64 `json:"confidence"`
	RangeDate  string  `json:"range_date"`
	CreatedAt  string  `json:"created_at"`
}

type UpdateAprioriRequest struct {
	IdApriori   int    `json:"id_apriori"`
	Code        string `json:"code"`
	Description string `json:"description"`
	Image       string `json:"image"`
}

type GenerateApriori struct {
	ItemSet     []string `json:"item_set"`
	Support     float32  `json:"support"`
	Iterate     int      `json:"iterate"`
	Transaction int      `json:"transaction"`
	Confidence  float32  `json:"confidence"`
	Discount    float32  `json:"discount"`
	Description string   `json:"description"`
	RangeDate   string   `json:"range_date"`
}

type GenerateAprioriRequest struct {
	MinimumSupport    float64 `json:"minimum_support"`
	MinimumConfidence float64 `json:"minimum_confidence"`
	MinimumDiscount   int     `json:"minimum_discount"`
	MaximumDiscount   int     `json:"maximum_discount"`
	StartDate         string  `json:"start_date"`
	EndDate           string  `json:"end_date"`
}

type GenerateCreateAprioriRequest struct {
	ItemSet     []string `json:"item_set"`
	Support     float64  `json:"support"`
	Iterate     int      `json:"iterate"`
	Transaction int      `json:"transaction"`
	Confidence  float64  `json:"confidence"`
	Discount    float64  `json:"discount"`
	Description string   `json:"description"`
	RangeDate   string   `json:"range_date"`
}
