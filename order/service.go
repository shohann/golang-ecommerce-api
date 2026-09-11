package order

import (
	"fmt"

	"github.com/shohann/golang-ecommerce-api/apperr"
	"github.com/shohann/golang-ecommerce-api/cart"
	"github.com/shohann/golang-ecommerce-api/config"
	"github.com/shohann/golang-ecommerce-api/domain"
	orderHandler "github.com/shohann/golang-ecommerce-api/rest/handlers/order"
)

type service struct {
	cartItemRepo cart.CartItemRepo
	orderRepo    OrderRepo
	publisher    OrderEventPublisher
	cnf          *config.Config
}

type Service interface {
	orderHandler.Service
}

func NewService(
	cartItemRepo cart.CartItemRepo,
	orderRepo OrderRepo,
	publisher OrderEventPublisher,
	cnf *config.Config,
) Service {
	return &service{
		cartItemRepo: cartItemRepo,
		orderRepo:    orderRepo,
		publisher:    publisher,
		cnf:          cnf,
	}
}

func (svc *service) OrderCheckOut(userId int64) (*domain.Order, error) {
	userCartItems, err := svc.cartItemRepo.GetCartByUserId(userId)

	if err != nil {
		return nil, apperr.WrapInternal("GetCartByUserId", err)
	}

	// Check if the cart is empty (len == 0 works for both nil and empty slices)
	if len(userCartItems) == 0 {
		return nil, apperr.NotFound("cart is empty")
	}

	totalCost := 0.0
	orderItems := make([]domain.OrderItem, 0, len(userCartItems))

	for _, cartItem := range userCartItems {
		if cartItem.Product == nil {
			return nil, apperr.Internal("checkout", fmt.Errorf("product missing for product_id: %d", cartItem.ProductID))
		}

		unitPrice := cartItem.Product.Price
		subtotal := float64(cartItem.Quantity) * unitPrice

		orderItem := domain.OrderItem{
			ProductID: cartItem.ProductID,
			Quantity:  cartItem.Quantity,
			UnitPrice: unitPrice,
			Subtotal:  subtotal,
			Product:   cartItem.Product,
		}

		orderItems = append(orderItems, orderItem)

		totalCost += subtotal
	}

	shippingAddress := "Test"

	orderData := domain.Order{
		UserID:          userId,
		TotalAmount:     totalCost,
		Status:          domain.OrderStatusPending,
		ShippingAddress: shippingAddress,
		Items:           orderItems,
	}

	createdOrder, err := svc.orderRepo.Create(orderData)
	if err != nil {
		return nil, apperr.WrapInternal("create_order", err)
	}

	err = svc.publisher.PublishOrderPlaced(OrderPlacedEvent{
		OrderID:     createdOrder.ID,
		UserID:      createdOrder.UserID,
		TotalAmount: createdOrder.TotalAmount,
		Status:      string(createdOrder.Status),
	})
	if err != nil {
		return nil, apperr.WrapInternal("publish_order_placed", err)
	}

	return createdOrder, nil
}
