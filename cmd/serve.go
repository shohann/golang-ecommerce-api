package cmd

import (
	"fmt"
	"os"

	"github.com/shohann/golang-ecommerce-api/config"
	"github.com/shohann/golang-ecommerce-api/infra/db"
	"github.com/shohann/golang-ecommerce-api/rest"
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

	// middlewares := middleware.NewMiddlewares(cnf)

	// repos
	// productRepo := repo.NewProductRepo(dbCon)
	// userRepo := repo.NewUserRepo(dbCon)

	// domains
	// usrSvc := user.NewService(userRepo)
	// prdctSvc := product.NewService(productRepo)

	// productHandler := prdctHandler.NewHandler(middlewares, prdctSvc)
	// userHandler := usrHandler.NewHandler(cnf, usrSvc)

	server := rest.NewServer(
		cnf,
		// productHandler,
		// userHandler,
	)
	server.Start()
}
