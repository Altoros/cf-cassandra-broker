package migrate

import (
	"fmt"

	"github.com/Altoros/cf-cassandra-broker/config"
	"github.com/gocql/gocql"
)

func Run(config *config.CassandraConfig) error {
	var err error

	session, err := connectToCassandra(config, false)
	if err != nil {
		return fmt.Errorf("error connecting to cassandra: %s", err.Error())
	}

	err = createKeyspace(session, config.Keyspace)
	if err != nil {
		return fmt.Errorf("error creating keyspace: %s", err.Error())
	}

	session.Close()

	session, err = connectToCassandra(config, true)
	if err != nil {
		return fmt.Errorf("error connecting to cassandra: %s", err.Error())
	}
	defer session.Close()

	err = createInstancesTable(session, config.Keyspace)
	if err != nil {
		return fmt.Errorf("error creating instances: %s", err.Error())
	}

	err = createBindingsTable(session, config.Keyspace)
	if err != nil {
		return fmt.Errorf("error creating bindings: %s", err.Error())
	}

	return nil
}

func connectToCassandra(config *config.CassandraConfig, useKeyspace bool) (*gocql.Session, error) {
	cluster := gocql.NewCluster(config.Nodes...)
	if useKeyspace {
		cluster.Keyspace = config.Keyspace
	}
	cluster.Consistency = gocql.Quorum
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: config.Username,
		Password: config.Password,
	}

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}

	return session, nil
}

func createKeyspace(session *gocql.Session, keyspace string) error {
	query := fmt.Sprintf(`
CREATE KEYSPACE IF NOT EXISTS %s
WITH replication = { 'class' : 'SimpleStrategy', 'replication_factor' : 3 }`, keyspace)
	err := session.Query(query).Consistency(gocql.Quorum).Exec()
	if err != nil {
		return err
	}
	return nil
}

func createInstancesTable(session *gocql.Session, keyspace string) error {
	createTableQuery := `
CREATE TABLE IF NOT EXISTS instances (
	id text PRIMARY KEY,
	keyspace_name text,
	created_at timestamp
)`

	err := session.Query(createTableQuery).Exec()
	if err != nil {
		return fmt.Errorf("failed to create table: %s", err.Error())
	}
	return nil
}

func createBindingsTable(session *gocql.Session, keyspace string) error {
	createTableQuery := `
CREATE TABLE IF NOT EXISTS bindings (
	id text PRIMARY KEY,
	instance_id text,
	app_guid text,
	username text,
	password text,
	created_at timestamp
)`
	err := session.Query(createTableQuery).Consistency(gocql.All).Exec()
	if err != nil {
		return fmt.Errorf("failed to create table: %s", err.Error())
	}

	return nil
}
