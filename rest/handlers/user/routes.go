package user

import (
	"net/http"

	middleware "github.com/shohann/golang-ecommerce-api/rest/middlewares"
)

func (h *Handler) RegisterRoutes(mux *http.ServeMux, manager *middleware.Manager) {
	mux.Handle("POST /users", manager.With(
		http.HandlerFunc(h.CreateUser),
	))

	mux.Handle("POST /users/login", manager.With(
		http.HandlerFunc(h.LoginUser),
	))

	mux.Handle("GET /users/me", manager.With(
		http.HandlerFunc(h.GetUserProfile),
		h.middlewares.Authenticate,
	))
}
