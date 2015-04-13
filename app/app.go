package app

import (
	"log"
	"os"

	"github.com/bmizerany/pat"
	"github.com/codegangsta/negroni"

	"github.com/Altoros/cf-cassandra-service-broker/config"
)

type AppContext struct {
	appConfig *config.Config
	negroni   *negroni.Negroni
}

func NewApp(appConfig *config.Config) (*AppContext, error) {
	var err error

	app := new(AppContext)
	app.appConfig = appConfig

	err = app.initApi()
	if err != nil {
		return nil, err
	}

	return app, nil
}

func (app *AppContext) Start() {
	log.Println("Starting broker on port", app.appConfig.PortStr())
	app.negroni.Run(":" + app.appConfig.PortStr())
}

func (app *AppContext) Stop() {
	log.Println("Stopping broker")
	os.Exit(0)
}

func (app *AppContext) initApi() error {
	mux := pat.New()
	app.negroni = negroni.Classic()
	app.negroni.UseHandler(mux)

	return nil
}
