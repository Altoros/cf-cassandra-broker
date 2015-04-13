package config

import (
	"io/ioutil"
	"strconv"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Username string        `yaml:"username"`
	Password string        `yaml:"password"`
	Port     uint16        `yaml:"port"`
	Catalog  CatalogConfig `yaml:"catalog"`
}

type CatalogConfig struct {
	Services []ServiceConfig `yaml:"services" json:"services"`
}

type ServiceConfig struct {
	Id          string                `yaml:"id"          json:"id"`
	Name        string                `yaml:"name"        json:"name"`
	Description string                `yaml:"description" json:"description"`
	Bindable    bool                  `yaml:"bindable"    json:"bindable"`
	Tags        []string              `yaml:"tags"        json:"tags"`
	Metadata    ServiceMetadataConfig `yaml:"metadata"    json:"metadata"`
	Plans       []PlanConfig          `yaml:"plans"       json:"plans"`
}

type ServiceMetadataConfig struct {
	DisplayName         string `yaml:"displayName"         json:"displayName"`
	DocumentationUrl    string `yaml:"documentationUrl"    json:"documentationUrl"`
	ImageUrl            string `yaml:"imageUrl"            json:"imageUrl"`
	LongDescription     string `yaml:"longDescription"     json:"longDescription"`
	ProviderDisplayName string `yaml:"providerDisplayName" json:"providerDisplayName"`
	SupportUrl          string `yaml:"supportUrl"          json:"supportUrl"`
}

type PlanConfig struct {
	Id          string             `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Metadata    PlanMetadataConfig `json:"metadata"`
}

type PlanMetadataConfig struct {
	DisplayName string           `yaml:"displayName" json:"displayName"`
	Costs       []PlanCostConfig `yaml:"costs"       json:"costs"`
}

type PlanCostConfig struct {
	Unit   string             `yaml:"unit"   json:"unit"`
	Amount map[string]float32 `yaml:"amount" json:"amount"`
}

var defaultConfig = Config{
	Port: 80,
}

func Default() *Config {
	config := defaultConfig

	// config.Process()

	return &config
}

func (c *Config) Initialize(configYAML []byte) error {
	return yaml.Unmarshal(configYAML, &c)
}

func (c *Config) PortStr() string {
	return strconv.Itoa(int(c.Port))
}

func InitFromFile(path string) (*Config, error) {
	var config *Config = Default()
	var err error

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = config.Initialize(bytes)
	if err != nil {
		return nil, err
	}

	// c.Process()

	return config, nil
}
