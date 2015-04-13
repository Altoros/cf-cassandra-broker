package api

import (
	"net/http"

	"github.com/Altoros/cf-cassandra-service-broker/config"

	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

type catalogController struct {
	render  *render.Render
	catalog *config.CatalogConfig
}

func NewCatalogController(catalog *config.CatalogConfig) *catalogController {
	catalogController := new(catalogController)
	catalogController.render = render.New()
	catalogController.catalog = catalog
	return catalogController
}

func (c *catalogController) AddRoutes(router *mux.Router) {
	router.HandleFunc("/catalog", c.catalogHandler).Methods("GET")
}

func (c *catalogController) catalogHandler(res http.ResponseWriter, req *http.Request) {
	c.render.JSON(res, http.StatusOK, c.catalog)
}
