package model

type Product struct {
	Name     string  `json:"name"`
	Category string  `json:"category"`
	Unit     float64 `json:"units"`
	Tax      float64 `json:"tax"`
	Price    float64 `json:"amount"`
	Vendorid int     `json:"vendor_id"`
	Status   bool    `json:"status"`
	Discount float64 `json:"discount"`
}
type ProductList struct {
	Name       string  `json:"name"`
	Category   string  `json:"category"`
	Unit       float64 `json:"units"`
	Tax        float64 `json:"tax"`
	Price      float64 `json:"amount"`
	VendorName string  `json:"vendorName"`
	Status     bool    `json:"status"`
	Discount   float64 `json:"discount"`
}
