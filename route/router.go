package route

import (
	"chi-demo/handler"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth/v5"
)

var tokenAuth *jwtauth.JWTAuth

func init() {
	tokenAuth = jwtauth.New("HS256", []byte("secret"), nil)
	_, tokenString, _ := tokenAuth.Encode(map[string]interface{}{"user_id": 123})
	fmt.Printf("DEBUG: a sample jwt is %s\n\n", tokenString)
}

func InitRouter(r *chi.Mux, productHandler handler.ProductHandler) {
	// Protected routes
	r.Group(func(r chi.Router) {
		// Seek, verify and validate JWT tokens
		r.Use(jwtauth.Verifier(tokenAuth))

		// Handle valid / invalid tokens
		r.Use(jwtauth.Authenticator)

		r.Get("/products/{id}", productHandler.GetOne())

		r.Get("/products", productHandler.GetProducts())

		r.Post("/product", productHandler.CreateProduct())

		r.Delete("/products/{id}", productHandler.DeleteProduct())

	})

	r.Group(func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("root."))
		})
	})

}
