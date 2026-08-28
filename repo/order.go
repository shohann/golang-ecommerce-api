package repo

import (
	"github.com/jmoiron/sqlx"
	"github.com/shohann/golang-ecommerce-api/apperr"
	"github.com/shohann/golang-ecommerce-api/domain"
	orderDomain "github.com/shohann/golang-ecommerce-api/order"
)

type OrderRepo interface {
	orderDomain.OrderRepo
}

type orderRepo struct {
	db *sqlx.DB
}

func NewOrderRepo(db *sqlx.DB) OrderRepo {
	return &orderRepo{db: db}
}

func (r *orderRepo) Create(order domain.Order) (*domain.Order, error) {
	createdOrder := order
	err := r.db.QueryRow(`
		INSERT INTO orders (user_id, total_amount, status, shipping_address)
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, total_amount, status, shipping_address, created_at, updated_at
	`, order.UserID, order.TotalAmount, order.Status, order.ShippingAddress).Scan(
		&createdOrder.ID,
		&createdOrder.UserID,
		&createdOrder.TotalAmount,
		&createdOrder.Status,
		&createdOrder.ShippingAddress,
		&createdOrder.CreatedAt,
		&createdOrder.UpdatedAt,
	)
	if err != nil {
		return nil, apperr.Internal("create_order", err)
	}

	createdOrder.Items = make([]domain.OrderItem, 0, len(order.Items))
	for _, item := range order.Items {
		createdItem := item
		err := r.db.QueryRow(`
			INSERT INTO order_items (order_id, product_id, quantity, unit_price, subtotal)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id, order_id, product_id, quantity, unit_price, subtotal
		`, createdOrder.ID, item.ProductID, item.Quantity, item.UnitPrice, item.Subtotal).Scan(
			&createdItem.ID,
			&createdItem.OrderID,
			&createdItem.ProductID,
			&createdItem.Quantity,
			&createdItem.UnitPrice,
			&createdItem.Subtotal,
		)
		if err != nil {
			return nil, apperr.Internal("create_order_item", err)
		}

		createdOrder.Items = append(createdOrder.Items, createdItem)
	}

	return &createdOrder, nil
}
