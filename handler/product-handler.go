package handler

import (
	"chi-demo/model"
	"chi-demo/service"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

type ProductHandler struct {
	ProductService service.ProductService
}

type HandlerErr struct {
	Code        int
	Description string
}

func (e HandlerErr) Error() string {
	return e.Description
}

func ErrHandler(handlerFunc func(w http.ResponseWriter, r *http.Request) error) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := handlerFunc(w, r); err != nil {
			herr, ok := err.(HandlerErr)
			if ok {
				w.WriteHeader(herr.Code)

				json.NewEncoder(w).Encode(model.Response{
					Code:        herr.Code,
					Description: herr.Description,
				})

				log.Printf("error %s\n", err.Error())

				return
			}

			w.WriteHeader(http.StatusInternalServerError)

			json.NewEncoder(w).Encode(model.Response{
				Code:        http.StatusInternalServerError,
				Description: "Internal Server Error",
			})

			log.Printf("error %s\n", err.Error())
		}
	}
}

func (productHandler ProductHandler) GetOne(w http.ResponseWriter, r *http.Request) error {
	idParam := chi.URLParam(r, "id")

	// convert id to int
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return HandlerErr{
			Code:        http.StatusBadRequest,
			Description: "Cannot convert id to integer",
		}
	}

	// check if id > 0
	if id < 0 {
		return HandlerErr{
			Code:        http.StatusBadRequest,
			Description: "Invalid id",
		}
	}

	product, err := productHandler.ProductService.GetOne(r.Context(), int64(id))
	if err != nil {
		return err
	}

	json.NewEncoder(w).Encode(product)
	return nil
}

func (productHandler ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) error {
	products, err := productHandler.ProductService.GetAll(r.Context())
	if err != nil {
		return err
	}
	json.NewEncoder(w).Encode(products)
	return nil
}

func (productHandler ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) error {
	var inputProduct model.Product
	if err := json.NewDecoder(r.Body).Decode(&inputProduct); err != nil {
		return HandlerErr{
			Code:        http.StatusBadRequest,
			Description: "Invalid product",
		}
	}

	// check if fields exist
	if inputProduct.Name == "" || inputProduct.Price == 0 {
		return HandlerErr{
			Code:        http.StatusBadRequest,
			Description: "Missing field",
		}
	}

	// check if price > 0
	if inputProduct.Price < 0 {
		return HandlerErr{
			Code:        http.StatusBadRequest,
			Description: "Invalid price",
		}
	}

	err := productHandler.ProductService.Create(r.Context(), inputProduct)
	if err != nil {
		// log.Printf("Error when create product: %s", err.Error())
		return err
	}

	json.NewEncoder(w).Encode(model.Response{
		Code:        http.StatusOK,
		Description: "Product created",
	})

	return nil
}

func (productHandler ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) error {
	idParam := chi.URLParam(r, "id")

	// convert id to int
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return HandlerErr{
			Code:        http.StatusBadRequest,
			Description: "Cannot convert id to integer",
		}
	}

	// check if id > 0
	if id < 0 {
		return HandlerErr{
			Code:        http.StatusBadRequest,
			Description: "Invalid id",
		}
	}

	err = productHandler.ProductService.Delete(r.Context(), int64(id))
	if err != nil {
		return err
	}

	json.NewEncoder(w).Encode(model.Response{
		Code:        http.StatusOK,
		Description: "Product deleted",
	})
	return nil
}
