package dto

type ItemMeta struct {
	ProductID int64
	Qty       int64
	Price     int64
}

type PublicCheckout struct {
	Items      []*ItemMeta `json:"items"`
	SessionKey string      `header:"session-id"`
}
