package config

type CassandraConfig struct {
	Nodes      []string `yaml:"nodes"`
	CqlPort    uint16   `yaml:"cql_port"`
	ThriftPort uint16   `yaml:"thrift_port"`
	Keyspace   string   `yaml:"keyspace"`
	Username   string   `yaml:"username"`
	Password   string   `yaml:"password"`
}

var defaultCassandraConfig = CassandraConfig{
	CqlPort:    9042,
	ThriftPort: 9160,
}
