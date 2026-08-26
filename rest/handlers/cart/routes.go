package cart

import (
	"net/http"

	middleware "github.com/shohann/golang-ecommerce-api/rest/middlewares"
)

func (h *Handler) RegisterRoutes(mux *http.ServeMux, manager *middleware.Manager) {
	mux.Handle("POST /cart-items", manager.With(
		http.HandlerFunc(h.AddCartItem),
		h.middlewares.Authenticate,
	))

	mux.Handle("GET /cart-items", manager.With(
		http.HandlerFunc(h.GetUserCart),
		h.middlewares.Authenticate,
	))

	mux.Handle("PUT /cart-items/{id}", manager.With(
		http.HandlerFunc(h.UpdateCartItem),
		h.middlewares.Authenticate,
	))

	mux.Handle("DELETE /cart-items/{id}", manager.With(
		http.HandlerFunc(h.DeleteCartItem),
		h.middlewares.Authenticate,
	))
}
