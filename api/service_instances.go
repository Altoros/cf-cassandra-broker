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
	router.HandleFunc("/service_instances/{instance_id}", c.create).Methods("PUT")
	router.HandleFunc("/service_instances/{instance_id}", c.delete).Methods("DELETE")
}

func (c *serviceInstancesController) create(res http.ResponseWriter, req *http.Request) {
	var err error

	vars := mux.Vars(req)
	instanceId := vars["instance_id"]

	instanceExist, err := c.cassandraService.IsInstanceExist(instanceId)
	if err != nil {
		log.Println("Create instance:", err.Error())
		renderer.JSON(res, http.StatusInternalServerError, emptyResponse)
		return
	}

	if instanceExist {
		renderer.JSON(res, http.StatusConflict, ApiError("instance %s already exist", instanceId))
	} else {
		err = c.cassandraService.CreateInstance(instanceId)
		if err != nil {
			log.Println("Create instance:", err.Error())
			renderer.JSON(res, http.StatusInternalServerError, emptyResponse)
			return
		}
		renderer.JSON(res, http.StatusCreated, emptyResponse)
	}
}

func (c *serviceInstancesController) delete(res http.ResponseWriter, req *http.Request) {
	var err error

	vars := mux.Vars(req)
	instanceId := vars["instance_id"]

	instanceExist, err := c.cassandraService.IsInstanceExist(instanceId)
	if err != nil {
		log.Println("Delete instance:", err.Error())
		renderer.JSON(res, http.StatusInternalServerError, emptyResponse)
		return
	}

	if instanceExist {
		err = c.cassandraService.DeleteInstance(instanceId)
		if err != nil {
			log.Println("Delete instance:", err.Error())
			renderer.JSON(res, http.StatusInternalServerError, emptyResponse)
			return
		}
		renderer.JSON(res, http.StatusOK, emptyResponse)
	} else {
		renderer.JSON(res, http.StatusGone, ApiError("instance %s does not exist", instanceId))
	}
}
