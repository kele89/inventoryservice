package product

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/kele89/inventoryservice/cors"
)

const productBasePath = "products"

func SetUpRoutes(apiBasePath string) {
	productsHandler := http.HandlerFunc(handleProducts)
	productHandler := http.HandlerFunc(handleProduct)
	http.Handle(fmt.Sprintf("%s/%s", apiBasePath, productBasePath), cors.Middleware(productsHandler))
	http.Handle(fmt.Sprintf("%s/%s/", apiBasePath, productBasePath), cors.Middleware(productHandler))
}

func handleProducts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		productList := getProductList()
		productJSON, err := json.Marshal(productList)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(productJSON)
	case http.MethodPost:
		var newProduct Product
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(bodyBytes, &newProduct)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if newProduct.ProductId != 0 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, err = addOrUpdateProduct(newProduct)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusCreated)
		return
	case http.MethodOptions:
		return
	}
}

func handleProduct(w http.ResponseWriter, r *http.Request) {
	urlPathSegments := strings.Split(r.URL.Path, "products/")
	productId, err := strconv.Atoi(urlPathSegments[len(urlPathSegments)-1])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	product := getProduct(productId)
	if product == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodGet:
		productJSON, err := json.Marshal(product)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(productJSON)
	case http.MethodPut:
		var updatedProduct Product
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(bodyBytes, &updatedProduct)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if updatedProduct.ProductId != productId {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		addOrUpdateProduct(updatedProduct)
		w.WriteHeader(http.StatusOK)
		return
	case http.MethodDelete:
		removeProduct(productId)
	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
