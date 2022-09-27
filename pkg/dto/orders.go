package dto

type ItemMeta struct {
	ProductID int64 `json:"product_id"`
	Qty       int64 `json:"qty"`
	Price     int64 `json:"price"`
}

type PublicCheckout struct {
	Items      []*ItemMeta `json:"items"`
	SessionKey string      `header:"session-id"`
}
