package product

import (
	"net/http"

	middleware "github.com/shohann/golang-ecommerce-api/rest/middlewares"
)

func (h *Handler) RegisterRoutes(mux *http.ServeMux, manager *middleware.Manager) {
	mux.Handle("POST /products", manager.With(
		http.HandlerFunc(h.CreateProduct),
		h.middlewares.Authenticate,
	))

	mux.Handle("GET /products", manager.With(http.HandlerFunc(h.ListProducts)))

	mux.Handle("DELETE /products/{id}", manager.With(
		http.HandlerFunc(h.DeleteProduct),
		h.middlewares.Authenticate,
	))
}
