package apriori

type CreateAprioriRequest struct {
	Item       string  `json:"item"`
	Discount   float32 `json:"discount"`
	Support    float32 `json:"support"`
	Confidence float32 `json:"confidence"`
	RangeDate  string  `json:"range_date"`
	CreatedAt  string  `json:"created_at"`
}

type UpdateAprioriRequest struct {
	IdApriori   int64
	Code        string
	Description string
	Image       string
}

type GenerateAprioriRequest struct {
	MinimumSupport    float32 `json:"minimum_support" binding:"required,max=100"`
	MinimumConfidence float32 `json:"minimum_confidence" binding:"required,max=100"`
	MinimumDiscount   int32   `json:"minimum_discount" binding:"required"`
	MaximumDiscount   int32   `json:"maximum_discount" binding:"required,gtefield=MinimumDiscount"`
	StartDate         string  `json:"start_date" binding:"required"`
	EndDate           string  `json:"end_date" binding:"required"`
}
