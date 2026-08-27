package order

import (
	"net/http"

	middleware "github.com/shohann/golang-ecommerce-api/rest/middlewares"
)

func (h *Handler) RegisterRoutes(mux *http.ServeMux, manager *middleware.Manager) {
	mux.Handle("POST /orders/checkout", manager.With(
		http.HandlerFunc(h.OrderCheckOut),
		h.middlewares.Authenticate,
	))
}
