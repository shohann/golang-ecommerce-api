package cmd

import (
	"fmt"
	"os"

	"github.com/shohann/golang-ecommerce-api/config"
	"github.com/shohann/golang-ecommerce-api/infra/db"
	"github.com/shohann/golang-ecommerce-api/repo"
	"github.com/shohann/golang-ecommerce-api/rest"
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
	// productRepo := repo.NewProductRepo(dbCon)
	userRepo := repo.NewUserRepo(dbCon)

	// domains
	usrSvc := user.NewService(userRepo, cnf)
	// prdctSvc := product.NewService(productRepo)

	// productHandler := prdctHandler.NewHandler(middlewares, prdctSvc)
	// userHandler := NewHandler(cnf, usrSvc)

	usrHandler := userHandler.NewHandler(cnf, usrSvc, middlewares)

	server := rest.NewServer(
		cnf,
		usrHandler,
	)
	server.Start()
}
