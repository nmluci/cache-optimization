package dto

type PublicProduct struct {
	Name        string `json:"name"`
	Category    string `json:"category"`
	Description string `json:"description"`
	UnitPrice   string `json:"unit_price"`
	Qty         uint64 `json:"qty"`
}
