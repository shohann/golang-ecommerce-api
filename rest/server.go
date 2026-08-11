package rest

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/shohann/golang-ecommerce-api/config"
	"github.com/shohann/golang-ecommerce-api/rest/handlers/user"
	middleware "github.com/shohann/golang-ecommerce-api/rest/middlewares"
)

type Server struct {
	cnf *config.Config
	// productHandler *product.Handler
	userHandler *user.Handler
}

func NewServer(
	cnf *config.Config,
	// productHandler *product.Handler,
	userHandler *user.Handler,
) *Server {
	return &Server{
		cnf:         cnf,
		userHandler: userHandler,
		// productHandler: productHandler,
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

	// server.productHandler.RegisterRoutes(mux, manager)
	server.userHandler.RegisterRoutes(mux, manager)

	addr := ":" + strconv.Itoa(server.cnf.HttpPort)
	fmt.Println("Server running on port", addr)
	err := http.ListenAndServe(addr, warappedMux)
	if err != nil {
		fmt.Println("Error starting the server", err)
		os.Exit(1)
	}
}
