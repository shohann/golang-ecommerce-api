package cmd

import (
	"fmt"
	"os"

	"github.com/shohann/golang-ecommerce-api/cart"
	"github.com/shohann/golang-ecommerce-api/category"
	"github.com/shohann/golang-ecommerce-api/config"
	"github.com/shohann/golang-ecommerce-api/infra/db"
	"github.com/shohann/golang-ecommerce-api/order"
	"github.com/shohann/golang-ecommerce-api/product"
	"github.com/shohann/golang-ecommerce-api/repo"
	"github.com/shohann/golang-ecommerce-api/rest"
	cartHandler "github.com/shohann/golang-ecommerce-api/rest/handlers/cart"
	categoryHandler "github.com/shohann/golang-ecommerce-api/rest/handlers/category"
	orderHandler "github.com/shohann/golang-ecommerce-api/rest/handlers/order"
	productsHandler "github.com/shohann/golang-ecommerce-api/rest/handlers/product"
	userHandler "github.com/shohann/golang-ecommerce-api/rest/handlers/user"
	middleware "github.com/shohann/golang-ecommerce-api/rest/middlewares"
	"github.com/shohann/golang-ecommerce-api/user"
)

func Serve() {
	cnf := config.GetConfig()

	dbCon, err := db.NewConnection(cnf.DB)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = db.MigrateDB(dbCon, "migrations")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	middlewares := middleware.NewMiddlewares(cnf)

	// repos
	userRepo := repo.NewUserRepo(dbCon)
	categoryRepo := repo.NewCategoryRepo(dbCon)
	productRepo := repo.NewProductRepo(dbCon)
	cartItemRepo := repo.NewCartItemRepo(dbCon)

	// domains
	usrSvc := user.NewService(userRepo, cnf)
	categorySvc := category.NewService(categoryRepo)
	prdctSvc := product.NewService(productRepo, cnf)
	cartItemSvc := cart.NewService(cartItemRepo, cnf)
	orderSvc := order.NewService(cartItemRepo, cnf)

	// handlers
	productHandler := productsHandler.NewHandler(cnf, middlewares, prdctSvc)
	usrHandler := userHandler.NewHandler(cnf, usrSvc, middlewares)
	catHandler := categoryHandler.NewHandler(cnf, categorySvc, middlewares)
	crtHandler := cartHandler.NewHandler(cnf, middlewares, cartItemSvc)
	orderHandler := orderHandler.NewHandler(cnf, middlewares, orderSvc)

	server := rest.NewServer(
		cnf,
		productHandler,
		usrHandler,
		catHandler,
		crtHandler,
		orderHandler,
	)
	server.Start()
}
