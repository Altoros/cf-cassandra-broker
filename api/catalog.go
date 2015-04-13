package api

import (
	"net/http"

	"github.com/Altoros/cf-cassandra-service-broker/config"

	"github.com/bmizerany/pat"
	"github.com/unrolled/render"
)

type catalogController struct {
	render  *render.Render
	catalog *config.CatalogConfig
}

func NewCatalog(catalog *config.CatalogConfig) *catalogController {
	catalogController := new(catalogController)
	catalogController.render = render.New()
	catalogController.catalog = catalog
	return catalogController
}

func (c *catalogController) AddRoutes(mux *pat.PatternServeMux) {
	mux.Get("/catalog", http.HandlerFunc(c.catalogHandler))
}

func (c *catalogController) catalogHandler(res http.ResponseWriter, req *http.Request) {
	c.render.JSON(res, http.StatusOK, c.catalog)
}
