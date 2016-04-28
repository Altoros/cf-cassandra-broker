package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/cloudfoundry-community/types-cf"
	"github.com/codegangsta/negroni"
	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"

	"github.com/Altoros/cf-cassandra-broker/config"
)

var (
	renderer      = render.New(render.Options{IndentJSON: true})
	emptyResponse = make(map[string]string)
)

type ApiHandler struct {
	Handler *negroni.Negroni
	Config  *config.Config
	Service ServiceProvider
}

func New(appConfig *config.Config, session *gocql.Session) http.Handler {
	apiHandler := new(ApiHandler)
	apiHandler.Config = appConfig

	apiLogger := NewLogger()
	panicRecovery := negroni.NewRecovery()
	panicRecovery.PrintStack = false
	panicRecovery.Logger = apiLogger.Logger
	apiHandler.Handler = negroni.New(apiLogger, panicRecovery)
	apiHandler.Service = &cassandraService{session: session}

	apiHandler.DefineRoutes()

	return apiHandler
}

func (a *ApiHandler) DefineRoutes() {
	router := mux.NewRouter()

	router.HandleFunc("/v2/catalog", a.ShowCatalog).Methods("GET")
	router.HandleFunc("/v2/service_instances/{instance_id}", a.CreateServiceInstance).Methods("PUT")
	router.HandleFunc("/v2/service_instances/{instance_id}", a.DeleteServiceInstance).Methods("DELETE")
	router.HandleFunc("/v2/service_instances/{instance_id}/service_bindings/{binding_id}", a.CreateServiceBinding).Methods("PUT")
	router.HandleFunc("/v2/service_instances/{instance_id}/service_bindings/{binding_id}", a.DeleteServiceBinding).Methods("DELETE")

	a.Handler.UseHandler(router)
}

func (a *ApiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.Handler.ServeHTTP(w, r)
}

func writeError(w http.ResponseWriter, err *cf.ServiceProviderError) {
	if err.Code < 500 {
		renderer.JSON(w, err.Code, cf.BrokerError{err.String()})
	}
}

func (a *ApiHandler) ShowCatalog(w http.ResponseWriter, r *http.Request) {
	renderer.JSON(w, http.StatusOK, a.Config.Catalog)
}

func (a *ApiHandler) CreateServiceInstance(w http.ResponseWriter, r *http.Request) {
	serviceCreationRequest := new(cf.ServiceCreationRequest)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	json.Unmarshal(body, serviceCreationRequest)

	serviceCreationRequest.InstanceID = mux.Vars(r)["instance_id"]

	serviceError := a.Service.CreateService(serviceCreationRequest)
	if serviceError == nil {
		renderer.JSON(w, http.StatusCreated, emptyResponse)
	} else {
		writeError(w, serviceError)
	}
}

func (a *ApiHandler) DeleteServiceInstance(w http.ResponseWriter, r *http.Request) {
	instanceId := mux.Vars(r)["instance_id"]
	serviceError := a.Service.DeleteService(instanceId)
	if serviceError == nil {
		renderer.JSON(w, http.StatusOK, emptyResponse)
	} else {
		writeError(w, serviceError)
	}
}

func (a *ApiHandler) CreateServiceBinding(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	serviceBindingRequest := new(cf.ServiceBindingRequest)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	json.Unmarshal(body, serviceBindingRequest)

	serviceBindingRequest.InstanceID = vars["instance_id"]
	serviceBindingRequest.BindingID = vars["binding_id"]

	serviceBindingResponse, serviceError := a.Service.BindService(serviceBindingRequest)

	if serviceError == nil {
		creds := &serviceBindingResponse.Credentials

		creds.Nodes = a.Config.Cassandra.Nodes
		creds.CqlPort = a.Config.Cassandra.CqlPort
		creds.ThriftPort = a.Config.Cassandra.ThriftPort

		renderer.JSON(w, http.StatusCreated, serviceBindingResponse)
	} else {
		writeError(w, serviceError)
	}
}

func (a *ApiHandler) DeleteServiceBinding(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	serviceError := a.Service.UnbindService(vars["instance_id"], vars["binding_id"])
	if serviceError == nil {
		renderer.JSON(w, http.StatusOK, emptyResponse)
	} else {
		writeError(w, serviceError)
	}
}
