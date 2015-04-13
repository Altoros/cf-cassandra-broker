package api_test

import (
	. "github.com/Altoros/cf-cassandra-service-broker/api"
	"github.com/Altoros/cf-cassandra-service-broker/config"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"net/http"
	"net/http/httptest"

	"github.com/gorilla/mux"
)

var _ = Describe("Catalog", func() {
	var request *http.Request
	var recorder *httptest.ResponseRecorder
	var router *mux.Router

	BeforeEach(func() {
		router = mux.NewRouter()
		recorder = httptest.NewRecorder()
	})

	Describe("GET /catalog", func() {
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
			catalog := config.CatalogConfig{
				Services: []config.ServiceConfig{service},
			}
			catalogController := NewCatalogController(&catalog)
			catalogController.AddRoutes(router)
			request, _ = http.NewRequest("GET", "/catalog", nil)
			router.ServeHTTP(recorder, request)
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
})
