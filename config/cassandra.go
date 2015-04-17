package config

import "strconv"

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

func (c *CassandraConfig) PortString() string {
	return strconv.Itoa(int(c.Port))
}
