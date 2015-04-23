# CF Cassandra Broker

CF Cassandra Broker provides Cassandra databases as a Cloud Foundry service. This broker demonstrates the v2 services API between cloud controllers and service brokers. This API is not to be confused with the cloud controller API; for more info see [http://docs.cloudfoundry.org/services/api.html].

The management tasks that the broker performs are as follows:

* Provisioning of database instances (create)
* Creation of credentials (bind)
* Removal of credentials (unbind)
* Unprovisioning of database instances (delete)

## Prerequesites

* Enabled password authentication for cassandra cluster, see [http://docs.datastax.com/en/cassandra/1.2/cassandra/security/security_config_native_authenticate_t.html]
* Existing superuser

## Testing

To run all specs: `ginkgo -r`

## Installation

go get http://github.com/altoros/cf-cassandra-broker/cmd/cf-cassandra-broker
go get http://github.com/altoros/cf-cassandra-broker/cmd/cf-cassandra-broker-migrate

## Usage

Configure the config file for your environment. See `config.yml.example` for example.

Run migrate tool to prepare broker administrative keyspace:

```
cf-cassandra-broker-migrate -c <path to config file>
```

Start the Cassandra Service Broker:

```
cf-cassandra-broker -c <path to config file>
```

Add the broker to Cloud Foundry as described by [the service broker documentation](http://docs.cloudfoundry.org/services/managing-service-brokers.html).
