package product

type ReqCreateProduct struct {
	CategoryID  int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	ImageURL    string  `json:"image_url"`
	IsActive    bool    `json:"is_active"`
}
