package model

import (
	"net/url"
)

type Product struct {
	Name        string  `json:"name"`
	Category    string  `json:"category"`
	Unit        float64 `json:"units"`
	Tax         float64 `json:"tax"`
	Price       float64 `json:"amount"`
	Vendorid    string  `json:"vendor_id"`
	Status      bool    `json:"status"`
	Discount    float64 `json:"discount"`
	Description string  `json:"description"`
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
type ProductListUsers struct {
	Name       string  `json:"name"`
	Category   string  `json:"category"`
	Unit       float64 `json:"units"`
	Tax        float64 `json:"tax"`
	Price      float64 `json:"amount"`
	VendorName string  `json:"vendorName"`
	Status     bool    `json:"status"`
	Discount   float64 `json:"discount"`
	Pid        string  `json:"pid"`
}

func (u *Product) Valid() url.Values {
	err := url.Values{}

	if u.Name == "" {
		err.Add("name", "Name is required")
	}
	if u.Category == "" {
		err.Add("category", "Category is required")
	}
	if u.Unit <= 0 {
		err.Add("unit", "Unit must be greater than zero")
	}
	if u.Tax < 0 {
		err.Add("tax", "Tax cannot be negative")
	}
	if u.Price <= 0 {

		err.Add("price", "Price must be greater than zero")
	}
	if u.Discount < 0 {
		err.Add("discount", "Discount cannot be negative")
	}
	if u.Description == "" {
		err.Add("description", "Description is required")
	}
	if u.Vendorid != "" {
		err.Add("irrelevant value", "Dont enter vendor_id no irrevalant values")

	}

	return err
}
