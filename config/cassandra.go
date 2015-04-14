package config

type CassandraConfig struct {
	Nodes    []string `yaml:"nodes"`
	Keyspace string   `yaml:"keyspace"`
	Username string   `yaml:"username"`
	Password string   `yaml:"password"`
}
