package api_test

import (
	"github.com/Altoros/cf-cassandra-service-broker/api"
	"github.com/Altoros/cf-cassandra-service-broker/config"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cloudfoundry-community/types-cf"
	"github.com/codegangsta/negroni"

	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
)

type mockCassandraService struct {
	InstanceExist bool
	BindingExist  bool
}

func (s *mockCassandraService) CreateService(r *cf.ServiceCreationRequest) *cf.ServiceProviderError {
	if s.InstanceExist {
		return cf.NewServiceProviderError(cf.ErrorInstanceExists, errors.New(r.InstanceID))
	}

	return nil
}

func (s *mockCassandraService) DeleteService(instanceID string) *cf.ServiceProviderError {
	if !s.InstanceExist {
		return cf.NewServiceProviderError(cf.ErrorInstanceNotFound, errors.New(instanceID))
	}

	return nil
}

func (s *mockCassandraService) BindService(r *cf.ServiceBindingRequest) (*api.ServiceBindingResponse, *cf.ServiceProviderError) {
	if !s.InstanceExist {
		return nil, cf.NewServiceProviderError(cf.ErrorInstanceNotFound, errors.New(r.InstanceID))
	}

	if s.BindingExist {
		return nil, cf.NewServiceProviderError(cf.ErrorInstanceExists, errors.New(r.BindingID))
	}

	response := &api.ServiceBindingResponse{
		Credentials: api.ServiceCredentials{
			Username: "username",
			Password: "password",
			Keyspace: "keyspace",
			Vhost:    "keyspace",
		},
	}
	return response, nil
}

func (s *mockCassandraService) UnbindService(instanceID, bindingID string) *cf.ServiceProviderError {
	if !s.InstanceExist {
		return cf.NewServiceProviderError(cf.ErrorInstanceNotFound, errors.New(instanceID))
	}

	if !s.BindingExist {
		return cf.NewServiceProviderError(cf.ErrorInstanceNotFound, errors.New(bindingID))
	}

	return nil
}

var _ = Describe("API", func() {
	var request *http.Request
	var recorder *httptest.ResponseRecorder
	var apiInstance api.ApiHandler
	var cassandraService *mockCassandraService

	BeforeEach(func() {
		cassandraService = &mockCassandraService{}
		apiInstance = api.ApiHandler{
			Handler: negroni.New(),
			Service: cassandraService,
			Config:  &config.Config{},
		}
		apiInstance.DefineRoutes()
		recorder = httptest.NewRecorder()
	})

	Describe("GET /v2/catalog", func() {
		BeforeEach(func() {
			service := config.ServiceConfig{
				Id:          "service id",
				Name:        "service name",
				Description: "service description",
				Bindable:    true,
				Tags:        []string{"foo", "bar"},
				Metadata: config.ServiceMetadataConfig{
					DisplayName:         "service name",
					DocumentationUrl:    "http://example.com",
					ImageUrl:            "http://example.com/logo.png",
					LongDescription:     "long description",
					ProviderDisplayName: "provider display name",
					SupportUrl:          "http://example.com",
				},
				Plans: []config.PlanConfig{
					config.PlanConfig{
						Id:          "plan id",
						Name:        "plan name",
						Description: "plan description",
						Metadata: config.PlanMetadataConfig{
							DisplayName: "keyspace",
							Costs: []config.PlanCostConfig{
								config.PlanCostConfig{
									Unit: "MONTHLY",
									Amount: map[string]float32{
										"usd": float32(0.0),
										"eur": float32(0.0),
									},
								},
							},
						},
					},
				},
			}

			apiInstance.Config.Catalog = config.CatalogConfig{
				Services: []config.ServiceConfig{service},
			}

			request, _ = http.NewRequest("GET", "/v2/catalog", nil)
			apiInstance.ServeHTTP(recorder, request)
		})

		It("returns a status code of 200", func() {
			Ω(recorder.Code).To(Equal(200))
		})

		It("returns application/json content type", func() {
			Ω(recorder.Header()["Content-Type"]).To(Equal([]string{"application/json; charset=UTF-8"}))
		})

		It("returns proper json", func() {
			Ω(recorder.Body).To(MatchJSON(`
{
	"services": [
		{
			"id":          "service id",
			"name":        "service name",
			"description": "service description",
			"bindable":    true,
			"tags":        ["foo", "bar"],
			"metadata": {
				"displayName":         "service name",
				"documentationUrl":    "http://example.com",
				"imageUrl":            "http://example.com/logo.png",
				"longDescription":     "long description",
				"providerDisplayName": "provider display name",
				"supportUrl":          "http://example.com"
			},
			"plans": [
				{
					"id":          "plan id",
					"name":        "plan name",
					"description": "plan description",
					"metadata":    {
						"displayName": "keyspace",
						"costs":       [
							{
								"unit":   "MONTHLY",
								"amount": {
									"usd": 0.0,
									"eur": 0.0
								}
							}
						]
					}
				}
			]
		}
	]
}
`))
		})
	})

	Describe("PUT /v2/service_instances/:instance_id", func() {
		Context("Instance does not exist", func() {
			BeforeEach(func() {
				cassandraService.InstanceExist = false
				request, _ = http.NewRequest("PUT", "/v2/service_instances/foobar", strings.NewReader("{}"))
				apiInstance.ServeHTTP(recorder, request)
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
				request, _ = http.NewRequest("PUT", "/v2/service_instances/foobar", strings.NewReader("{}"))
				apiInstance.ServeHTTP(recorder, request)
			})

			It("returns a status code of 409", func() {
				Ω(recorder.Code).To(Equal(409))
			})

			It("returns application/json content type", func() {
				Ω(recorder.Header()["Content-Type"]).To(Equal([]string{"application/json; charset=UTF-8"}))
			})

			It("returns empty json", func() {
				Ω(recorder.Body).To(MatchJSON(`{"description":  "Error: 409 (ErrorInstanceExists) - foobar"}`))
			})
		})
	})

	Describe("DELETE /service_instances/:instance_id", func() {
		Context("Instance does not exist", func() {
			BeforeEach(func() {
				cassandraService.InstanceExist = false
				request, _ = http.NewRequest("DELETE", "/v2/service_instances/foobar", nil)
				apiInstance.ServeHTTP(recorder, request)
			})

			It("returns a status code of 410", func() {
				Ω(recorder.Code).To(Equal(410))
			})

			It("returns empty json", func() {
				Ω(recorder.Body).To(MatchJSON(`{"description": "Error: 410 (ErrorInstanceNotFound) - foobar"}`))
			})
		})

		Context("Instance exists", func() {
			BeforeEach(func() {
				cassandraService.InstanceExist = true
				request, _ = http.NewRequest("DELETE", "/v2/service_instances/foobar", nil)
				apiInstance.ServeHTTP(recorder, request)
			})

			It("returns a status code of 200", func() {
				Ω(recorder.Code).To(Equal(200))
			})

			It("returns application/json content type", func() {
				Ω(recorder.Header()["Content-Type"]).To(Equal([]string{"application/json; charset=UTF-8"}))
			})

			It("returns empty json", func() {
				Ω(recorder.Body).To(MatchJSON("{}"))
			})
		})
	})

	Describe("PUT /v2/service_instances/:instance_id/service_bindings/:binding_id", func() {
		Context("Instance does not exist", func() {
			BeforeEach(func() {
				cassandraService.InstanceExist = false
				body := strings.NewReader("{}")
				request, _ = http.NewRequest("PUT", "/v2/service_instances/foo/service_bindings/bar", body)
				apiInstance.ServeHTTP(recorder, request)
			})

			It("returns a status code of 410", func() {
				Ω(recorder.Code).To(Equal(410))
			})

			It("returns application/json content type", func() {
				Ω(recorder.Header()["Content-Type"]).To(Equal([]string{"application/json; charset=UTF-8"}))
			})

			It("returns empty json", func() {
				Ω(recorder.Body).To(MatchJSON(`{"description": "Error: 410 (ErrorInstanceNotFound) - foo"}`))
			})
		})

		Context("Instance exists", func() {
			BeforeEach(func() {
				cassandraService.InstanceExist = true
			})

			Context("Binding exists", func() {
				BeforeEach(func() {
					cassandraService.BindingExist = true
					request, _ = http.NewRequest("PUT", "/v2/service_instances/foo/service_bindings/bar", strings.NewReader("{}"))
					apiInstance.ServeHTTP(recorder, request)
				})

				It("returns a status code of 409", func() {
					Ω(recorder.Code).To(Equal(409))
				})

				It("returns application/json content type", func() {
					Ω(recorder.Header()["Content-Type"]).To(Equal([]string{"application/json; charset=UTF-8"}))
				})

				It("returns json with error", func() {
					Ω(recorder.Body).To(MatchJSON(`{"description":  "Error: 409 (ErrorInstanceExists) - bar"}`))
				})
			})

			Context("Binding does not exists", func() {
				BeforeEach(func() {
					apiInstance.Config.Cassandra = config.CassandraConfig{
						Nodes: []string{"host1", "host2"},
						Port:  123,
					}
					cassandraService.InstanceExist = true
					request, _ = http.NewRequest("PUT", "/v2/service_instances/foo/service_bindings/bar", strings.NewReader("{}"))
					apiInstance.ServeHTTP(recorder, request)
				})

				It("returns a status code of 201", func() {
					Ω(recorder.Code).To(Equal(201))
				})

				It("returns application/json content type", func() {
					Ω(recorder.Header()["Content-Type"]).To(Equal([]string{"application/json; charset=UTF-8"}))
				})

				It("returns json with credentials", func() {
					Ω(recorder.Body).To(MatchJSON(`
{
	"credentials": {
		"username": "username",
		"password": "password",
		"nodes": ["host1", "host2"],
		"port": "123",
		"keyspace": "keyspace",
		"vhost": "keyspace"
	}
}`))
				})
			})
		})
	})

	Describe("DELETE /v2/service_instances/:instance_id/service_bindings/:binding_id", func() {
		Context("Instance does not exist", func() {
			BeforeEach(func() {
				cassandraService.InstanceExist = false
				request, _ = http.NewRequest("DELETE", "/v2/service_instances/foo/service_bindings/bar", nil)
				apiInstance.ServeHTTP(recorder, request)
			})

			It("returns a status code of 410", func() {
				Ω(recorder.Code).To(Equal(410))
			})

			It("returns application/json content type", func() {
				Ω(recorder.Header()["Content-Type"]).To(Equal([]string{"application/json; charset=UTF-8"}))
			})
		})

		Context("Instance exists", func() {
			BeforeEach(func() {
				cassandraService.InstanceExist = true
			})

			Context("Binding exists", func() {
				BeforeEach(func() {
					cassandraService.BindingExist = true
					request, _ = http.NewRequest("DELETE", "/v2/service_instances/foo/service_bindings/bar", strings.NewReader("{}"))
					apiInstance.ServeHTTP(recorder, request)
				})

				It("returns a status code of 200", func() {
					Ω(recorder.Code).To(Equal(200))
				})

				It("returns application/json content type", func() {
					Ω(recorder.Header()["Content-Type"]).To(Equal([]string{"application/json; charset=UTF-8"}))
				})

				It("returns empty json", func() {
					Ω(recorder.Body).To(MatchJSON(`{}`))
				})
			})

			Context("Binding does not exists", func() {
				BeforeEach(func() {
					cassandraService.BindingExist = false
					request, _ = http.NewRequest("DELETE", "/v2/service_instances/foo/service_bindings/bar", strings.NewReader("{}"))
					apiInstance.ServeHTTP(recorder, request)
				})

				It("returns a status code of 410", func() {
					Ω(recorder.Code).To(Equal(410))
				})

				It("returns application/json content type", func() {
					Ω(recorder.Header()["Content-Type"]).To(Equal([]string{"application/json; charset=UTF-8"}))
				})
			})
		})
	})
})
