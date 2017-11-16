package broker

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gocql/gocql"
	"github.com/goji/httpauth"

	"github.com/Altoros/cf-cassandra-broker/api"
	"github.com/Altoros/cf-cassandra-broker/config"
)

type AppContext struct {
	config           *config.Config
	serveMux         *http.ServeMux
	cassandraSession *gocql.Session
}

func New(appConfig *config.Config) (*AppContext, error) {
	app := new(AppContext)
	app.config = appConfig
	session, err := newCassandraSession(&appConfig.Cassandra)
	if err != nil {
		return nil, fmt.Errorf("can't start cassandra session: %s", err)
	}
	app.cassandraSession = session

	app.serveMux = http.NewServeMux()
	apiAuthHandler := httpauth.SimpleBasicAuth(appConfig.Username, appConfig.Password)
	api := api.New(app.config, app.cassandraSession)
	app.serveMux.Handle("/v2/", apiAuthHandler(api))

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

func newCassandraSession(cfg *config.CassandraConfig) (*gocql.Session, error) {
	cluster := gocql.NewCluster(cfg.Nodes...)
	cluster.Keyspace = cfg.Keyspace
	cluster.Timeout = 1 * time.Minute
	cluster.NumConns = 1
	if len(cfg.Nodes) == 1 {
		cluster.Consistency = gocql.One
	} else {
		cluster.Consistency = gocql.All
	}
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: cfg.Username,
		Password: cfg.Password,
	}

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}
	return session, nil
}
