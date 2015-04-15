package broker

import (
	"log"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"

	"github.com/Altoros/cf-cassandra-service-broker/api"
	"github.com/Altoros/cf-cassandra-service-broker/common"
	"github.com/Altoros/cf-cassandra-service-broker/config"
)

type AppContext struct {
	config           *config.Config
	negroni          *negroni.Negroni
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
	log.Println("Starting broker on port", app.config.PortStr())
	app.negroni.Run(":" + app.config.PortStr())
}

func (app *AppContext) Stop() {
	log.Println("Stopping broker")
	app.cassandraService.Stop()
}

func (app *AppContext) initCassandra() error {
	var err error

	cassandraCfg := app.config.Cassandra
	app.cassandraService, err = common.NewCassandraService(cassandraCfg.Nodes, cassandraCfg.Keyspace,
		cassandraCfg.Username, cassandraCfg.Password)

	return err
}

func (app *AppContext) initApi() error {
	negroni := negroni.New(
		negroni.NewRecovery(),
		negroni.NewLogger(),
	)
	mainRouter := mux.NewRouter()
	negroni.UseHandler(mainRouter)
	app.negroni = negroni

	apiRouter := mainRouter.PathPrefix("/v2").Subrouter()

	catalogController := api.NewCatalogController(&app.config.Catalog)
	catalogController.AddRoutes(apiRouter)

	serviceInstancesController := api.NewServiceInstancesController(app.cassandraService)
	serviceInstancesController.AddRoutes(apiRouter)

	return nil
}
