package api

import (
	"log"
	"net/http"

	"github.com/Altoros/cf-cassandra-service-broker/common"

	"github.com/gorilla/mux"
)

type serviceInstancesController struct {
	cassandraService common.CassandraService
}

func NewServiceInstancesController(cassandraService common.CassandraService) *serviceInstancesController {
	serviceInstancesController := new(serviceInstancesController)
	serviceInstancesController.cassandraService = cassandraService
	return serviceInstancesController
}

func (c *serviceInstancesController) AddRoutes(router *mux.Router) {
	router.HandleFunc("/service_instances/{instance_id}", c.createServiceInstance).Methods("PUT")
}

func (c *serviceInstancesController) createServiceInstance(res http.ResponseWriter, req *http.Request) {
	var err error

	vars := mux.Vars(req)
	instanceId := vars["instance_id"]

	instanceExist, err := c.cassandraService.IsInstanceExist(instanceId)
	if err != nil {
		log.Println(err.Error())
		renderer.JSON(res, http.StatusInternalServerError, emptyResponse)
	} else {
		if instanceExist {
			renderer.JSON(res, http.StatusConflict, emptyResponse)
		} else {
			err = c.cassandraService.CreateInstance(instanceId)
			if err != nil {
				log.Println(err.Error())
				renderer.JSON(res, http.StatusInternalServerError, emptyResponse)
			} else {
				renderer.JSON(res, http.StatusCreated, emptyResponse)
			}
		}
	}
}
