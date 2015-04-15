package broker

import (
	"log"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"

	"github.com/Altoros/cf-cassandra-service-broker/api"
	"github.com/Altoros/cf-cassandra-service-broker/common"
	"github.com/Altoros/cf-cassandra-service-broker/config"
)

type AppContext struct {
	config           *config.Config
	apiHandler       *negroni.Negroni
	cassandraService common.CassandraService
}

func New(appConfig *config.Config) (*AppContext, error) {
	var err error

	app := new(AppContext)
	app.config = appConfig

	err = app.initCassandra()
	if err != nil {
		return nil, err
	}

	err = app.initApi()
	if err != nil {
		return nil, err
	}

	return app, nil
}

func (app *AppContext) Start() {
	port := app.config.PortStr()
	log.Println("Starting broker on port", port)
	http.ListenAndServe(":"+port, app.apiHandler)
}

func (app *AppContext) Stop() {
	log.Println("Stopping broker")
	app.cassandraService.Stop()
}

func (app *AppContext) initCassandra() error {
	var err error

	cfg := app.config.Cassandra
	app.cassandraService, err = common.NewCassandraService(cfg.Nodes, cfg.Keyspace,
		cfg.Username, cfg.Password)

	return err
}

func (app *AppContext) initApi() error {
	apiHandler := negroni.New(
		negroni.NewRecovery(),
		api.NewLogger(),
	)
	mainRouter := mux.NewRouter()
	apiHandler.UseHandler(mainRouter)
	app.apiHandler = apiHandler

	apiRouter := mainRouter.PathPrefix("/v2").Subrouter()

	catalogController := api.NewCatalogController(&app.config.Catalog)
	catalogController.AddRoutes(apiRouter)

	serviceInstancesController := api.NewServiceInstancesController(app.cassandraService)
	serviceInstancesController.AddRoutes(apiRouter)

	return nil
}
