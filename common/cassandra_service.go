package common

import (
	"time"

	"github.com/gocql/gocql"
)

type CassandraService interface {
	Stop()
	IsInstanceExist(instanceId string) (bool, error)
	CreateInstance(instanceId string) error
	DeleteInstance(instanceId string) error
}

type cassandraService struct {
	cluster *gocql.ClusterConfig
	session *gocql.Session
}

func NewCassandraService(hosts []string, keyspace string, username string, password string) (CassandraService, error) {
	var err error

	cs := new(cassandraService)

	cs.cluster = gocql.NewCluster(hosts...)
	cs.cluster.Keyspace = keyspace
	cs.cluster.Consistency = gocql.One
	cs.cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: username,
		Password: password,
	}

	cs.session, err = cs.cluster.CreateSession()
	if err != nil {
		return nil, err
	}

	return cs, nil
}

func (cs *cassandraService) Stop() {
	cs.session.Close()
}

func (cs *cassandraService) IsInstanceExist(instanceId string) (bool, error) {
	var records int
	err := cs.session.Query("SELECT COUNT(*) FROM instances WHERE id = ?", instanceId).Scan(&records)
	if err != nil {
		return false, err
	}

	return records > 0, nil
}

func (cs *cassandraService) CreateInstance(instanceId string) error {
	var err error

	err = cs.session.Query("CREATE KEYSPACE " + instanceId +
		" WITH replication = {'class': 'SimpleStrategy', 'replication_factor' : 3};").Exec()
	if err != nil {
		return err
	}

	err = cs.session.Query("INSERT INTO instances(id, created_at) VALUES(?, ?)", instanceId, time.Now()).Exec()
	if err != nil {
		return err
	}

	return nil
}

func (cs *cassandraService) DeleteInstance(instanceId string) error {
	var err error

	err = cs.session.Query("DELETE FROM instances WHERE id=?", instanceId).Exec()
	if err != nil {
		return err
	}

	err = cs.session.Query("DROP KEYSPACE " + instanceId).Exec()
	if err != nil {
		return err
	}

	return nil
}
