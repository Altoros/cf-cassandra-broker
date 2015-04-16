package broker

import (
	"log"
	"net/http"

	"github.com/gocql/gocql"

	"github.com/Altoros/cf-cassandra-service-broker/api"
	"github.com/Altoros/cf-cassandra-service-broker/config"
)

type AppContext struct {
	config           *config.Config
	serveMux         *http.ServeMux
	cassandraSession *gocql.Session
}

func New(appConfig *config.Config) (*AppContext, error) {
	app := new(AppContext)
	app.config = appConfig

	cluster := gocql.NewCluster(appConfig.Cassandra.Nodes...)
	cluster.Keyspace = appConfig.Cassandra.Keyspace
	cluster.Consistency = gocql.One
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: appConfig.Cassandra.Username,
		Password: appConfig.Cassandra.Password,
	}

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}
	app.cassandraSession = session

	app.serveMux = http.NewServeMux()

	api := api.New(&app.config.Catalog, app.cassandraSession)
	app.serveMux.Handle("/v2/", api)

	return app, nil
}

func (app *AppContext) Start() {
	port := app.config.PortStr()
	log.Println("Start broker on port", port)
	http.ListenAndServe(":"+port, app.serveMux)
}

func (app *AppContext) Stop() {
	log.Println("Stop broker")
	app.cassandraSession.Close()
}
