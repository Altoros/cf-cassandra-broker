package fakes

import (
	"errors"
)

type FakeCassandraService struct {
	InstanceExist bool
}

func (cs *FakeCassandraService) Stop() {
}

func (cs *FakeCassandraService) IsInstanceExist(instanceId string) (bool, error) {
	return cs.InstanceExist, nil
}

func (cs *FakeCassandraService) CreateInstance(instanceId string) error {
	if cs.InstanceExist {
		return errors.New("instance already exist")
	}
	return nil
}

func (cs *FakeCassandraService) DeleteInstance(instanceId string) error {
	if !cs.InstanceExist {
		return errors.New("instance does not exist")
	}
	return nil
}
