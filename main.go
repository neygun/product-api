package main

import (
	"chi-demo/db"
	"chi-demo/handler"
	"chi-demo/repository"
	"chi-demo/route"
	"chi-demo/service"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	_ "github.com/lib/pq"

	"chi-demo/log"
)

func router(productHandler handler.ProductHandler) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)

	route.InitRouter(r, productHandler)

	return r
}

func main() {
	// admin := model.User{
	// 	Username: "nguyen",
	// 	Password: "nguyen",
	// }

	logger := log.GetLogger()

	db, err := db.Init()
	if err != nil {
		panic(err)
	}

	productRepo := repository.New(db)
	productService := service.New(productRepo)
	productHandler := handler.New(productService)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3333"
	}

	logger.Printf("Running on port %s\n", port)
	http.ListenAndServe(":"+port, router(productHandler))
}
