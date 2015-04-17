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
		return err
	}

	err = createKeyspace(session, config.Keyspace)
	if err != nil {
		return err
	}

	session.Close()

	session, err = connectToCassandra(config, true)
	if err != nil {
		return err
	}
	defer session.Close()

	err = createInstancesTable(session, config.Keyspace)
	if err != nil {
		return fmt.Errorf("error creating instances table: %s", err.Error())
	}

	err = createBindingsTable(session, config.Keyspace)
	if err != nil {
		return fmt.Errorf("error creating bindings table: %s", err.Error())
	}

	return nil
}

func connectToCassandra(config *config.CassandraConfig, useKeyspace bool) (*gocql.Session, error) {
	cluster := gocql.NewCluster(config.Nodes...)
	if useKeyspace {
		cluster.Keyspace = config.Keyspace
	}
	// cluster.Consistency = gocql.One
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
	err := session.Query(query).Exec()
	if err != nil {
		return err
	}
	return nil
}

func createInstancesTable(session *gocql.Session, keyspace string) error {
	exist, err := isTableExist(session, keyspace, "instances")
	if err != nil {
		return err
	}

	if exist {
		return nil
	}

	createTableQuery := `
CREATE TABLE instances (
	id text PRIMARY KEY,
	keyspace_name text,
	created_at timestamp
)`

	err = session.Query(createTableQuery).Exec()
	if err != nil {
		return err
	}
	return nil
}

func createBindingsTable(session *gocql.Session, keyspace string) error {
	var err error

	exist, err := isTableExist(session, keyspace, "bindings")
	if err != nil {
		return err
	}

	if exist {
		return nil
	}

	createTableQuery := `
CREATE TABLE bindings (
	id text PRIMARY KEY,
	instance_id text,
	app_guid text,
	username text,
	password text,
	created_at timestamp
)`
	err = session.Query(createTableQuery).Exec()
	if err != nil {
		return err
	}

	err = session.Query("CREATE INDEX ON bindings (instance_id)").Exec()
	if err != nil {
		return err
	}

	return nil
}

func isTableExist(session *gocql.Session, keyspace string, table string) (bool, error) {
	var count int
	query := `
SELECT COUNT(*)
FROM system.schema_columnfamilies
WHERE keyspace_name=? AND columnfamily_name=?`
	err := session.Query(query, keyspace, table).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
