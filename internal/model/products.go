package model

// type Product struct {
// 	ProductID       uint64  `json:"product_id"`
// 	Name            string  `json:"name"`
// 	SupplierID      uint64  `json:"supplier_id"`
// 	CategoryID      uint64  `json:"category_id"`
// 	QuantityPerUnit string  `json:"qtyPerUnit"`
// 	UnitPrice       float32 `json:"unit_price"`
// 	Stocks          uint64  `json:"stocks"`
// 	OnOrder         uint64  `json:"on_order"`
// 	ReorderLevel    uint64  `json:"reoder_level"`
// 	IsDiscontinued  bool    `json:"is_discontinued"`
// }

// type Supplier struct {
// 	SupplierID   uint64 `json:"supplier_id"`
// 	CompanyName  string `json:"company_name"`
// 	ContactName  string `json:"contact_name"`
// 	ContactTitle string `json:"contact_title"`
// 	Address      string `json:"address"`
// 	City         string `json:"city"`
// 	Region       string `json:"region"`
// 	PostalCode   string `json:"postal_code"`
// 	Country      string `json:"country"`
// 	Phone        string `json:"phone"`
// 	Fax          string `json:"fax"`
// 	HomePage     string `json:"home_page"`
// }

// type Category struct {
// 	CategoryID  uint64 `json:"category_id"`
// 	Name        string `json:"name"`
// 	Description string `json:"description"`
// 	Picture     []byte
// }

type Product struct {
	ID          uint64 `json:"product_id"`
	Name        string `json:"name"`
	Category    string `json:"category"`
	Description string `json:"description"`
	UnitPrice   string `json:"unit_price"`
	Qty         uint64 `json:"qty"`
}
