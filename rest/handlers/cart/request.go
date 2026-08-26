package cart

type ReqAddCartItem struct {
	ProductID int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
}

type ReqUpdateCartItem struct {
	Quantity int `json:"quantity"`
}
