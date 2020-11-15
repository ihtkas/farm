package store

import (
	"context"

	"github.com/gocql/gocql"
)

// InitStoreInCassandra initializes data model for seller service in the given cassandra server.
func InitStoreInCassandra(ctx context.Context, clusterHosts []string, keyspace string) error {
	cluster := gocql.NewCluster(clusterHosts...)
	cluster.Keyspace = keyspace
	session, err := cluster.CreateSession()
	if err != nil {
		return err
	}

	q := session.Query(createOrderCasandraQuery)
	q = q.WithContext(ctx)
	err = q.Exec()
	if err != nil {
		return err
	}

	q = session.Query(createUserOrderCasandraQuery)
	q = q.WithContext(ctx)
	err = q.Exec()
	if err != nil {
		return err
	}
	return nil
}
