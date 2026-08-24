package domain

import "time"

type Cart struct {
	ID        int64      `db:"id"`
	UserID    int64      `db:"user_id"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	Items     []CartItem `db:"-"` // Optional relationship field for application logic
}

type CartItem struct {
	ID        int64     `db:"id"`
	CartID    int64     `db:"cart_id"`
	ProductID int64     `db:"product_id"`
	Quantity  int       `db:"quantity"`
	CreatedAt time.Time `db:"created_at"`
	Product   *Product  `db:"-"` // Optional relationship field for populated queries
}
