package rest

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/shohann/golang-ecommerce-api/config"
	"github.com/shohann/golang-ecommerce-api/rest/handlers/cart"
	"github.com/shohann/golang-ecommerce-api/rest/handlers/category"
	"github.com/shohann/golang-ecommerce-api/rest/handlers/product"
	"github.com/shohann/golang-ecommerce-api/rest/handlers/user"
	middleware "github.com/shohann/golang-ecommerce-api/rest/middlewares"
)

type Server struct {
	cnf             *config.Config
	productHandler  *product.Handler
	userHandler     *user.Handler
	categoryHandler *category.Handler
	cartHandler     *cart.Handler
}

func NewServer(
	cnf *config.Config,
	productHandler *product.Handler,
	userHandler *user.Handler,
	categoryHandler *category.Handler,
	cartHandler *cart.Handler,
) *Server {
	return &Server{
		cnf:             cnf,
		userHandler:     userHandler,
		productHandler:  productHandler,
		categoryHandler: categoryHandler,
		cartHandler:     cartHandler,
	}
}

func (server *Server) Start() {
	manager := middleware.NewManager()
	manager.Use(
		middleware.Preflight,
		middleware.Cors,
		middleware.Logger,
	)

	mux := http.NewServeMux()
	warappedMux := manager.WrapMux(mux)

	server.productHandler.RegisterRoutes(mux, manager)
	server.userHandler.RegisterRoutes(mux, manager)
	server.categoryHandler.RegisterRoutes(mux, manager)

	addr := ":" + strconv.Itoa(server.cnf.HttpPort)
	fmt.Println("Server running on port", addr)
	err := http.ListenAndServe(addr, warappedMux)
	if err != nil {
		fmt.Println("Error starting the server", err)
		os.Exit(1)
	}
}
