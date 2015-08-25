package api

import (
	"errors"
	"fmt"
	"time"

	"github.com/Altoros/cf-cassandra-broker/random"
	"github.com/cloudfoundry-community/types-cf"
	"github.com/gocql/gocql"
)

type ServiceProvider interface {
	// CreateService creates a service instance for specific plan
	CreateService(r *cf.ServiceCreationRequest) *cf.ServiceProviderError

	// DeleteService deletes previously created service instance
	DeleteService(instanceID string) *cf.ServiceProviderError

	// BindService binds to specified service instance and
	// Returns credentials necessary to establish connection to that service
	BindService(r *cf.ServiceBindingRequest) (*ServiceBindingResponse, *cf.ServiceProviderError)

	// UnbindService removes previously created binding
	UnbindService(instanceID, bindingID string) *cf.ServiceProviderError
}

type ServiceBindingResponse struct {
	Credentials ServiceCredentials `json:"credentials"`
}

type ServiceCredentials struct {
	Nodes      []string `json:"nodes"`
	CqlPort    uint16   `json:"cql_port"`
	ThriftPort uint16   `json:"thrift_port"`
	Username   string   `json:"username"`
	Password   string   `json:"password"`
	Keyspace   string   `json:"keyspace"`
}

type cassandraService struct {
	session *gocql.Session
}

// CreateService creates a service instance for specific plan
func (service *cassandraService) CreateService(r *cf.ServiceCreationRequest) *cf.ServiceProviderError {
	var err error

	if service.isInstanceExist(r.InstanceID) {
		return cf.NewServiceProviderError(cf.ErrorInstanceExists, errors.New(r.InstanceID))
	}

	keyspace := "cf" + random.Hex(10)

	query := "CREATE KEYSPACE " + keyspace +
		" WITH replication = {'class': 'SimpleStrategy', 'replication_factor' : 3};"
	err = service.session.Query(query).Exec()
	if err != nil {
		panic(err.Error())
	}

	err = service.session.Query("INSERT INTO instances(id, keyspace_name, created_at) VALUES(?, ?, ?)",
		r.InstanceID, keyspace, time.Now()).Exec()
	if err != nil {
		panic(err.Error())
	}

	return nil
}

// DeleteService deletes previously created service instance
func (service *cassandraService) DeleteService(instanceID string) *cf.ServiceProviderError {
	var err error

	if !service.isInstanceExist(instanceID) {
		return cf.NewServiceProviderError(cf.ErrorInstanceNotFound, errors.New(instanceID))
	}

	keyspace, err := service.findKeyspaceNameByInstanceId(instanceID)
	if err != nil {
		panic(err.Error())
	}

	var bindingID, username string
	iter := service.session.Query("SELECT id, username FROM bindings WHERE instance_id=?", instanceID).Iter()
	for iter.Scan(&bindingID, &username) {
		err = service.deleteBinding(bindingID)
		if err = iter.Close(); err != nil {
			panic(err.Error())
		}
		err = service.dropUser(username)
		if err = iter.Close(); err != nil {
			panic(err.Error())
		}
	}
	if err = iter.Close(); err != nil {
		panic(err.Error())
	}

	err = service.session.Query("DELETE FROM instances WHERE id=?", instanceID).Exec()
	if err != nil {
		panic(err.Error())
	}

	err = service.dropKeyspaceIfExist(keyspace)
	if err != nil {
		panic(err.Error())
	}

	return nil
}

// BindService binds to specified service instance and
// Returns credentials necessary to establish connection to that service
func (service *cassandraService) BindService(r *cf.ServiceBindingRequest) (*ServiceBindingResponse, *cf.ServiceProviderError) {
	var err error
	var query string

	if !service.isInstanceExist(r.InstanceID) {
		return nil, cf.NewServiceProviderError(cf.ErrorInstanceNotFound, errors.New(r.InstanceID))
	}

	if service.isBindingExist(r.BindingID) {
		return nil, cf.NewServiceProviderError(cf.ErrorInstanceExists, errors.New(r.BindingID))
	}

	username := "cf-" + random.Hex(10)
	password := random.Hex(10)
	keyspace, err := service.findKeyspaceNameByInstanceId(r.InstanceID)
	if err != nil {
		panic(err.Error())
	}

	query = fmt.Sprintf("CREATE USER '%s' WITH PASSWORD '%s' NOSUPERUSER", username, password)
	err = service.session.Query(query).Exec()
	if err != nil {
		panic(err.Error())
	}

	query = fmt.Sprintf("GRANT ALL PERMISSIONS on KEYSPACE %s TO \"%s\"", keyspace, username)
	err = service.session.Query(query).Exec()
	if err != nil {
		panic(err.Error())
	}

	err = service.session.Query(`INSERT INTO
		bindings(id, instance_id, app_guid, username, password, created_at)
		VALUES(?, ?, ?, ?, ?, ?)`,
		r.BindingID, r.InstanceID, r.AppGUID, username, password, time.Now()).Exec()
	if err != nil {
		panic(err.Error())
	}

	response := &ServiceBindingResponse{
		Credentials: ServiceCredentials{
			Username: username,
			Password: password,
			Keyspace: keyspace,
		},
	}

	return response, nil
}

// UnbindService removes previously created binding
func (service *cassandraService) UnbindService(instanceID, bindingID string) *cf.ServiceProviderError {
	var err error

	if !service.isInstanceExist(instanceID) {
		return cf.NewServiceProviderError(cf.ErrorInstanceNotFound, errors.New(instanceID))
	}

	var username string
	query := "SELECT username FROM bindings WHERE id = ? AND instance_id = ?"
	err = service.session.Query(query, bindingID, instanceID).Scan(&username)
	if err != nil {
		if err == gocql.ErrNotFound {
			return cf.NewServiceProviderError(cf.ErrorInstanceNotFound, errors.New(bindingID))
		} else {
			panic(err.Error())
		}
	}

	err = service.dropUser(username)
	if err != nil {
		panic(err.Error())
	}

	err = service.deleteBinding(bindingID)
	if err != nil {
		panic(err.Error())
	}

	return nil
}

func (service *cassandraService) isInstanceExist(instanceID string) bool {
	var recordsCount int

	query := "SELECT COUNT(*) FROM instances WHERE id = ?"
	err := service.session.Query(query, instanceID).Scan(&recordsCount)
	if err != nil {
		panic(err.Error())
	}

	return recordsCount > 0
}

func (service *cassandraService) isBindingExist(bindingID string) bool {
	var recordsCount int

	query := "SELECT COUNT(*) FROM bindings WHERE id = ?"
	err := service.session.Query(query, bindingID).Scan(&recordsCount)
	if err != nil {
		panic(err.Error())
	}

	return recordsCount > 0
}

func (service *cassandraService) findKeyspaceNameByInstanceId(instanceID string) (string, error) {
	var keyspace string
	query := "SELECT keyspace_name FROM instances WHERE id = ?"
	err := service.session.Query(query, instanceID).Scan(&keyspace)
	if err != nil {
		return "", err
	}
	return keyspace, nil
}

func (service *cassandraService) dropUser(name string) error {
	query := fmt.Sprintf("DROP USER '%s'", name)
	err := service.session.Query(query).Exec()
	if err != nil {
		return err
	}
	return nil
}

func (service *cassandraService) dropKeyspaceIfExist(keyspace string) error {
	var query string
	var err error

	var count int
	query = "SELECT COUNT(*) FROM system.schema_keyspaces WHERE keyspace_name=?"
	err = service.session.Query(query, keyspace).Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		return nil
	}

	err = service.session.Query("DROP KEYSPACE " + keyspace).Exec()
	if err != nil {
		return err
	}

	return nil
}

func (service *cassandraService) deleteBinding(bindingID string) error {
	query := "DELETE FROM bindings WHERE id = ?"
	err := service.session.Query(query, bindingID).Exec()
	if err != nil {
		return err
	}
	return nil
}
