package indto

type ItemData struct {
	ProductID int64
	Qty       int64
	Price     int64
}

type OrderData struct {
	UserdataID int64
	Orderdata  int64
	Items      []*ItemData
}
