package category

import (
	"net/http"

	middleware "github.com/shohann/golang-ecommerce-api/rest/middlewares"
)

func (h *Handler) RegisterRoutes(mux *http.ServeMux, manager *middleware.Manager) {
	mux.Handle("POST /categories", manager.With(
		http.HandlerFunc(h.CreateCategory),
		h.middlewares.RequireAdmin,
		h.middlewares.Authenticate,
	))

	mux.Handle("GET /categories", manager.With(
		http.HandlerFunc(h.ListCategories),
		h.middlewares.RequireAdmin,
		h.middlewares.Authenticate,
	))

	mux.Handle("GET /categories/{id}", manager.With(
		http.HandlerFunc(h.GetCategory),
		h.middlewares.RequireAdmin,
		h.middlewares.Authenticate,
	))

	mux.Handle("PUT /categories/{id}", manager.With(
		http.HandlerFunc(h.UpdateCategory),
		h.middlewares.RequireAdmin,
		h.middlewares.Authenticate,
	))

	mux.Handle("DELETE /categories/{id}", manager.With(
		http.HandlerFunc(h.DeleteCategory),
		h.middlewares.RequireAdmin,
		h.middlewares.Authenticate,
	))
}
