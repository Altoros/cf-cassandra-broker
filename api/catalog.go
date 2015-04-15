package api

import (
	"net/http"

	"github.com/Altoros/cf-cassandra-service-broker/config"

	"github.com/gorilla/mux"
)

type catalogController struct {
	catalog *config.CatalogConfig
}

func NewCatalogController(catalog *config.CatalogConfig) *catalogController {
	catalogController := new(catalogController)
	catalogController.catalog = catalog
	return catalogController
}

func (c *catalogController) AddRoutes(router *mux.Router) {
	router.HandleFunc("/catalog", c.show).Methods("GET")
}

func (c *catalogController) show(res http.ResponseWriter, req *http.Request) {
	renderer.JSON(res, http.StatusOK, c.catalog)
}
