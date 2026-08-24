package repo

import (
	"github.com/jmoiron/sqlx"
	"github.com/shohann/golang-ecommerce-api/apperr"
	"github.com/shohann/golang-ecommerce-api/cart"
	"github.com/shohann/golang-ecommerce-api/domain"
)

type CartItemRepo interface {
	cart.CartItemRepo
}

type cartItemRepo struct {
	db *sqlx.DB
}

func NewCartItemRepo(db *sqlx.DB) CartItemRepo {
	return &cartItemRepo{
		db: db,
	}
}

func (r *cartItemRepo) AddItem(cartItem domain.CartItem) (*domain.CartItem, error) {
	query := `
        INSERT INTO cart_items (
            cart_id,
            product_id,
            quantity
        )
        VALUES ($1, $2, $3)
        ON CONFLICT (cart_id, product_id) 
        DO UPDATE SET 
            quantity = cart_items.quantity + EXCLUDED.quantity
        RETURNING
            id, cart_id, product_id, quantity, created_at
    `

	err := r.db.QueryRow(
		query,
		cartItem.CartID,
		cartItem.ProductID,
		cartItem.Quantity,
	).Scan(
		&cartItem.ID,
		&cartItem.CartID,
		&cartItem.ProductID,
		&cartItem.Quantity,
		&cartItem.CreatedAt,
	)

	if err != nil {
		return nil, apperr.Internal("add_cart_item_repo", err)
	}

	return &cartItem, nil
}
