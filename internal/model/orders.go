package model

import "time"

// type Order struct {
// 	OrderID        uint64     `json:"order_id"`
// 	CustomerID     uint64     `json:"customer_id"`
// 	EmployeeID     uint64     `json:"employee_id"`
// 	OrderDate      *time.Time `json:"order_date"`
// 	RequiredDate   *time.Time `json:"required_date"`
// 	ShippedDate    *time.Time `json:"shipped_date"`
// 	ShipVia        uint64     `json:"ship_via"`
// 	Freight        float32    `json:"freight"`
// 	ShipName       string     `json:"ship_name"`
// 	ShipAddress    string     `json:"ship_address"`
// 	ShipCity       string     `json:"ship_city"`
// 	ShipRegion     string     `json:"ship_region"`
// 	ShipPostalCode string     `json:"ship_postal"`
// 	ShipCountry    string     `json:"ship_country"`
// }

type Order struct {
	ID        uint64    `json:"order_id"`
	UserID    uint64    `json:"user_id"`
	OrderDate time.Time `json:"order_date"`
	PaidAt    time.Time `json:"paid_at"`
}

type OrderDetail struct {
	ID        uint64 `json:"orderdetail_id"`
	OrderID   uint64 `json:"order_id"`
	ProductID uint64 `json:"product_id"`
	Qty       uint64 `json:"qty"`
	Price     uint64 `json:"price"`
}
