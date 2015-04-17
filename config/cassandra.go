package config

type CassandraConfig struct {
	Nodes    []string `yaml:"nodes"`
	Port     uint16   `yaml:"port"`
	Keyspace string   `yaml:"keyspace"`
	Username string   `yaml:"username"`
	Password string   `yaml:"password"`
}

var defaultCassandraConfig = CassandraConfig{
	Port: 9042,
}
