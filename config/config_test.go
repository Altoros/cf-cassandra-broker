package config_test

import (
	. "github.com/Altoros/cf-cassandra-broker/config"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {
	var config *Config

	BeforeEach(func() {
		config = Default()
	})

	Describe("defaults", func() {
		It("sets default value for port", func() {
			Ω(config.Port).To(Equal(uint16(80)))
		})

		Context("Cassandra", func() {
			It("sets default value for cql port", func() {
				Ω(config.Cassandra.CqlPort).To(Equal(uint16(9042)))
			})
			It("sets default value for thrift port", func() {
				Ω(config.Cassandra.ThriftPort).To(Equal(uint16(9160)))
			})
		})
	})

	Describe("Initialize", func() {
		It("sets catalog config", func() {
			var b = []byte(`
catalog:
  services:
  - bindable: true
    name: cassandra
    description: cassandra
    id: service-id
    metadata:
      displayName: cassandra
      documentationUrl: http://example.com
      imageUrl: http://example.com/logo.png
      longDescription: cassandra
      providerDisplayName: cassandra
      supportUrl: http://example.com
    plans:
    - name: free
      description: plan desc
      id: plan-id
      metadata:
        costs:
        - amount:
            usd: 0.0
          unit: MONTHLY
        displayName: keyspace
    tags:
    - cassandra
    - nosql
`)
			config.Initialize(b)

			Ω(len(config.Catalog.Services)).To(Equal(1))
			// Services
			Ω(config.Catalog.Services[0].Id).To(Equal("service-id"))
			Ω(config.Catalog.Services[0].Name).To(Equal("cassandra"))
			Ω(config.Catalog.Services[0].Description).To(Equal("cassandra"))
			Ω(config.Catalog.Services[0].Bindable).To(BeTrue())
			Ω(len(config.Catalog.Services[0].Tags)).To(Equal(2))
			// Service metadata
			Ω(config.Catalog.Services[0].Metadata.DisplayName).To(Equal("cassandra"))
			Ω(config.Catalog.Services[0].Metadata.DocumentationUrl).To(Equal("http://example.com"))
			Ω(config.Catalog.Services[0].Metadata.ImageUrl).To(Equal("http://example.com/logo.png"))
			Ω(config.Catalog.Services[0].Metadata.LongDescription).To(Equal("cassandra"))
			Ω(config.Catalog.Services[0].Metadata.ProviderDisplayName).To(Equal("cassandra"))
			Ω(config.Catalog.Services[0].Metadata.SupportUrl).To(Equal("http://example.com"))
			// Plans
			Ω(len(config.Catalog.Services[0].Plans)).To(Equal(1))
			Ω(config.Catalog.Services[0].Plans[0].Id).To(Equal("plan-id"))
			Ω(config.Catalog.Services[0].Plans[0].Name).To(Equal("free"))
			Ω(config.Catalog.Services[0].Plans[0].Description).To(Equal("plan desc"))
			// Plan metadata
			Ω(config.Catalog.Services[0].Plans[0].Metadata.DisplayName).To(Equal("keyspace"))
			// Costs
			Ω(len(config.Catalog.Services[0].Plans[0].Metadata.Costs)).To(Equal(1))
			Ω(config.Catalog.Services[0].Plans[0].Metadata.Costs[0].Unit).To(Equal("MONTHLY"))
			Ω(config.Catalog.Services[0].Plans[0].Metadata.Costs[0].Amount["usd"]).To(Equal(float32(0.0)))

		})

		It("sets cassandra config", func() {
			var b = []byte(`
cassandra:
  nodes:
  - 127.0.0.1
  keyspace: broker
  username: username
  password: password
  cql_port: 123
  thrift_port: 456
`)

			config.Initialize(b)
			Ω(config.Cassandra.Nodes).To(Equal([]string{"127.0.0.1"}))
			Ω(config.Cassandra.Keyspace).To(Equal("broker"))
			Ω(config.Cassandra.Username).To(Equal("username"))
			Ω(config.Cassandra.Password).To(Equal("password"))
			Ω(config.Cassandra.CqlPort).To(Equal(uint16(123)))
			Ω(config.Cassandra.ThriftPort).To(Equal(uint16(456)))
		})

		It("sets username", func() {
			var b = []byte(`
username: user
`)
			config.Initialize(b)
			Ω(config.Username).To(Equal("user"))
		})

		It("sets password", func() {
			var b = []byte(`
password: password
`)
			config.Initialize(b)
			Ω(config.Password).To(Equal("password"))
		})

		It("sets port", func() {
			var b = []byte(`
port: 8080
`)
			config.Initialize(b)
			Ω(config.Port).To(Equal(uint16(8080)))
		})
	})

	Describe("PortStr", func() {
		It("returns port as string", func() {
			config.Port = 1234
			Ω(config.PortStr()).To(Equal("1234"))
		})
	})
})
