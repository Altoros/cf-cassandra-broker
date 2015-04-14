package api_test

import (
	. "github.com/Altoros/cf-cassandra-service-broker/api"
	"github.com/Altoros/cf-cassandra-service-broker/common/fakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"net/http"
	"net/http/httptest"

	"github.com/gorilla/mux"
)

var _ = Describe("Service instances", func() {
	var (
		request          *http.Request
		recorder         *httptest.ResponseRecorder
		router           *mux.Router
		cassandraService *fakes.FakeCassandraService
	)

	BeforeEach(func() {
		router = mux.NewRouter()
		recorder = httptest.NewRecorder()
		cassandraService = &fakes.FakeCassandraService{}
		serviceInstancesController := NewServiceInstancesController(cassandraService)
		serviceInstancesController.AddRoutes(router)
	})

	Describe("PUT /service_instances/:instance_id", func() {
		Context("Instance does not exist", func() {
			BeforeEach(func() {
				cassandraService.InstanceExist = false
				request, _ = http.NewRequest("PUT", "/service_instances/foobar", nil)
				router.ServeHTTP(recorder, request)
			})

			It("returns a status code of 201", func() {
				Ω(recorder.Code).To(Equal(201))
			})

			It("returns application/json content type", func() {
				Ω(recorder.Header()["Content-Type"]).To(Equal([]string{"application/json; charset=UTF-8"}))
			})

			It("returns empty json", func() {
				Ω(recorder.Body).To(MatchJSON("{}"))
			})
		})

		Context("Instance exists", func() {
			BeforeEach(func() {
				cassandraService.InstanceExist = true
				request, _ = http.NewRequest("PUT", "/service_instances/foobar", nil)
				router.ServeHTTP(recorder, request)
			})

			It("returns a status code of 409", func() {
				Ω(recorder.Code).To(Equal(409))
			})

			It("returns application/json content type", func() {
				Ω(recorder.Header()["Content-Type"]).To(Equal([]string{"application/json; charset=UTF-8"}))
			})

			It("returns empty json", func() {
				Ω(recorder.Body).To(MatchJSON("{}"))
			})
		})
	})
})
