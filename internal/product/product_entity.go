package product

import (
	"time"
)

type Product struct {
	ID          string 			`json:"id" db:"id"`
	Name        string 			`json:"name"`
	SKU         string 			`json:"sku"`
	Category    ProductCategory `json:"category"`
	ImageURL    string 			`json:"imageUrl" db:"image_url"`
	Notes       string 			`json:"notes"`
	Price       float64 		`json:"price"`
	Stock       int 			`json:"stock"`
	Location    string 			`json:"location"`
	IsAvailable bool 			`json:"isAvailable" db:"is_available"`
	CreatedAt   time.Time 		`json:"createdAt" db:"created_at"`
}

type ProductCategory string

const (
	Clothing    ProductCategory = "Clothing"
	Accessories ProductCategory = "Accessories"
	Footwear    ProductCategory = "Footwear"
	Beverages   ProductCategory = "Beverages"
)

type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required,min=1,max=30"`
	Sku         string  `json:"sku" binding:"required,min=1,max=30"`
	Category    string  `json:"category" binding:"required,oneof=Clothing Accessories Footwear Beverages"`
	ImageURL    string  `json:"imageUrl" binding:"required"`
	Notes       string  `json:"notes" binding:"required,min=1,max=200"`
	Price       float64 `json:"price" binding:"required,min=1"`
	Stock       int     `json:"stock" binding:"required,min=0,max=100000"`
	Location    string  `json:"location" binding:"required,min=1,max=200"`
	IsAvailable bool    `json:"isAvailable" binding:"required"`
}

type CreateProductResponse struct {
	ID        string `json:"id"`
	CreatedAt string `json:"createdAt"`
}

func FormatCreateProductResponse(product Product) CreateProductResponse {
	return CreateProductResponse{
		ID:        product.ID,
		CreatedAt: product.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

type QueryParams struct {
	ID			string 	`form:"id"`
	Limit       int    	`form:"limit"`
	Offset      int    	`form:"offset"`
	Name        string 	`form:"name"`
	IsAvailable string 	`form:"isAvailable" binding:"omitempty,oneof=true false"`
	Category	string	`form:"category" binding:"omitempty,oneof=Clothing Accessories Footwear Beverages"`
	Sku        	string 	`form:"sku"`
	Price 		string	`form:"price" binding:"omitempty,oneof=asc desc"`
	InStock 	string 	`form:"inStock" binding:"omitempty,oneof=true false"`
	CreatedAt 	string	`form:"createdAt" binding:"omitempty,oneof=asc desc"`
}

type ProductFilter struct {
	Limit    int             `form:"id default=5"`
	Offset   int             `form:"offset default=0"`
	Name     string          `form:"name" json:"name"`
	Category ProductCategory `form:"category"`
	Sku      string          `form:"sku"`
	Price    SortBy          `form:"price"`
	InStock  InStockEnum     `form:"inStock"`
}

type InStockEnum string

const (
	True  InStockEnum = "true"
	False InStockEnum = "false"
)

type SortBy string

const (
	PriceAsc  SortBy = "asc"
	PriceDesc SortBy = "desc"
)
