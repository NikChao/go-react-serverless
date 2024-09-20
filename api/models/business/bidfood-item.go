package businessmodel

type BidFoodItem struct {
	Name  string  `json:"name"`
	Size  string  `json:"size"`
	UOM   string  `json:"uom"`
	Price float64 `json:"price"`
}
