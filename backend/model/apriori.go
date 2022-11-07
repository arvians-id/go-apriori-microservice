package model

import "time"

type Apriori struct {
	IdApriori   int64     `json:"id_apriori"`
	Code        string    `json:"code"`
	Item        string    `json:"item"`
	Discount    float32   `json:"discount"`
	Support     float32   `json:"support"`
	Confidence  float32   `json:"confidence"`
	RangeDate   string    `json:"range_date"`
	IsActive    bool      `json:"is_active"`
	Description *string   `json:"description"`
	Mass        int32     `json:"mass"`
	Image       *string   `json:"image"`
	CreatedAt   time.Time `json:"created_at"`
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
