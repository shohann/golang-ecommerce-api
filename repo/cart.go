package repo

import (
	"database/sql"
	"fmt"

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

func (r *cartItemRepo) GetOrCreateCartID(userID int64) (int64, error) {
	query := `
		INSERT INTO carts (user_id)
		VALUES ($1)
		ON CONFLICT (user_id) DO UPDATE SET updated_at = NOW()
		RETURNING id
	`

	var cartID int64
	err := r.db.QueryRow(query, userID).Scan(&cartID)
	if err != nil {
		return 0, apperr.Internal("get_or_create_cart", err)
	}

	return cartID, nil
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

	fmt.Println(cartItem)

	if err != nil {
		return nil, apperr.Internal("add_cart_item_repo", err)
	}

	return &cartItem, nil
}

func (r *cartItemRepo) GetCartByUserId(userId int64) ([]domain.CartItem, error) {
	query := `
		SELECT 
			ci.id,
			ci.cart_id,
			ci.product_id,
			ci.quantity,
			ci.created_at,
			p.id,
			p.category_id,
			p.name,
			p.description,
			p.price,
			p.stock,
			p.image_url,
			p.is_active,
			p.created_at,
			p.updated_at
		FROM cart_items ci
		JOIN carts c ON ci.cart_id = c.id
		JOIN products p ON ci.product_id = p.id
		WHERE c.user_id = $1
	`

	rows, err := r.db.Query(query, userId)
	if err != nil {
		return nil, apperr.Internal("get_cart_by_user_id", err)
	}
	defer rows.Close()

	items := make([]domain.CartItem, 0)

	for rows.Next() {
		var item domain.CartItem
		var product domain.Product

		err := rows.Scan(
			&item.ID,
			&item.CartID,
			&item.ProductID,
			&item.Quantity,
			&item.CreatedAt,
			&product.ID,
			&product.CategoryID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.Stock,
			&product.ImageURL,
			&product.IsActive,
			&product.CreatedAt,
			&product.UpdatedAt,
		)

		if err != nil {
			return nil, apperr.Internal("scan_cart_item", err)
		}

		item.Product = &product
		items = append(items, item)

	}

	if err = rows.Err(); err != nil {
		return nil, apperr.Internal("iterate_cart_items", err)
	}

	return items, nil
}

func (r *cartItemRepo) UpdateQuantity(userID, itemID int64, quantity int) (*domain.CartItem, error) {
	query := `
		UPDATE cart_items ci
		SET quantity = $1
		FROM carts c
		WHERE ci.cart_id = c.id
			AND c.user_id = $2
			AND ci.id = $3
		RETURNING ci.id, ci.cart_id, ci.product_id, ci.quantity, ci.created_at
	`

	var item domain.CartItem
	err := r.db.QueryRow(query, quantity, userID, itemID).Scan(
		&item.ID,
		&item.CartID,
		&item.ProductID,
		&item.Quantity,
		&item.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apperr.NotFound("cart item not found")
		}
		return nil, apperr.Internal("update_cart_item_quantity", err)
	}

	return &item, nil
}

func (r *cartItemRepo) DeleteItem(userID, itemID int64) error {
	query := `
		DELETE FROM cart_items ci
		USING carts c
		WHERE ci.cart_id = c.id
			AND c.user_id = $1
			AND ci.id = $2
	`

	result, err := r.db.Exec(query, userID, itemID)
	if err != nil {
		return apperr.Internal("delete_cart_item", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return apperr.Internal("delete_cart_item_rows", err)
	}
	if rows == 0 {
		return apperr.NotFound("cart item not found")
	}

	return nil
}
