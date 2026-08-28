package order

import "github.com/shohann/golang-ecommerce-api/domain"

type OrderRepo interface {
	Create(order domain.Order) (*domain.Order, error)
}
