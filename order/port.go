package order

import "github.com/shohann/golang-ecommerce-api/domain"

type OrderRepo interface {
	Create(order domain.Order) (*domain.Order, error)
}

type OrderPlacedEvent struct {
	OrderID     int64   `json:"order_id"`
	UserID      int64   `json:"user_id"`
	TotalAmount float64 `json:"total_amount"`
	Status      string  `json:"status"`
}

type OrderEventPublisher interface {
	PublishOrderPlaced(event OrderPlacedEvent) error
}
