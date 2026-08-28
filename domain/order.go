package domain

import "time"

type UserRole string

const (
	UserRoleAdmin    UserRole = "admin"
	UserRoleCustomer UserRole = "customer"
)

type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusSucceeded PaymentStatus = "succeeded"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusRefunded  PaymentStatus = "refunded"
)

type RefundStatus string

const (
	RefundStatusPending   RefundStatus = "pending"
	RefundStatusSucceeded RefundStatus = "succeeded"
	RefundStatusFailed    RefundStatus = "failed"
)

type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusPaid       OrderStatus = "paid"
	OrderStatusProcessing OrderStatus = "processing"
	OrderStatusShipped    OrderStatus = "shipped"
	OrderStatusDelivered  OrderStatus = "delivered"
	OrderStatusCancelled  OrderStatus = "cancelled"
)

type Order struct {
	ID              int64       `db:"id" json:"id"`
	UserID          int64       `db:"user_id" json:"user_id"`
	TotalAmount     float64     `db:"total_amount" json:"total_amount"`
	Status          OrderStatus `db:"status" json:"status"`
	ShippingAddress string      `db:"shipping_address" json:"shipping_address"`
	CreatedAt       time.Time   `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time   `db:"updated_at" json:"updated_at"`
	Items           []OrderItem `db:"-" json:"items,omitempty"`
}

type OrderItem struct {
	ID        int64    `db:"id" json:"id"`
	OrderID   int64    `db:"order_id" json:"order_id"`
	ProductID int64    `db:"product_id" json:"product_id"`
	Quantity  int      `db:"quantity" json:"quantity"`
	UnitPrice float64  `db:"unit_price" json:"unit_price"`
	Subtotal  float64  `db:"subtotal" json:"subtotal"`
	Product   *Product `db:"-" json:"product,omitempty"`
}
