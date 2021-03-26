package main

import (
	"net/http"

	"github.com/kele89/inventoryservice/product"
)

const apiBasePath = "/api"

func main() {
	product.SetUpRoutes(apiBasePath)
	http.ListenAndServe(":5000", nil)
}
