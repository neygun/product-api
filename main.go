package main

import (
	"chi-demo/handler"
	"chi-demo/repository"
	"chi-demo/service"
	"database/sql"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	_ "github.com/lib/pq"
)

func main() {
	connectionStr := "postgres://postgres:postgres@localhost:5432/scc-pg?sslmode=disable"

	db, err := sql.Open("postgres", connectionStr)
	if err != nil {
		panic(err)
	}

	productRepo := repository.New(db)
	productService := service.ProductServiceImpl{
		ProductRepository: productRepo,
	}
	productHandler := handler.ProductHandler{
		ProductService: productService,
	}

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)

	r.Get("/products/{id}", handler.ErrHandler(productHandler.GetOne))

	r.Get("/products", handler.ErrHandler(productHandler.GetProducts))

	r.Post("/product", handler.ErrHandler(productHandler.CreateProduct))

	r.Delete("/products/{id}", handler.ErrHandler(productHandler.DeleteProduct))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("root."))
	})

	http.ListenAndServe(":3333", r)
}
