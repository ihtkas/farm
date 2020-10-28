package store

import (
	"context"

	"github.com/gocql/gocql"
	"github.com/golang/glog"
	pgpoolv4 "github.com/jackc/pgx/v4/pgxpool"
)

// InitStoreInPostgres initializes data model for seller service in the given postgres server.
func InitStoreInPostgres(ctx context.Context, username, password, host, port, dbname string) error {
	config := "user=" + username +
		" host=" + host +
		" port=" + port +
		" dbname=" + dbname

	if password != "" {
		config += " password=" + password
	}

	glog.Errorln(config)
	pool, err := pgpoolv4.Connect(ctx, config)
	if err != nil {
		return err
	}

	conn, err := pool.Acquire(ctx)
	if err != nil {
		return err
	}

	tx, err := conn.Begin(ctx)

	_, err = tx.Exec(ctx, `CREATE EXTENSION IF NOT EXISTS postgis`)

	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)

	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `CREATE TABLE IF NOT EXISTS product 
	    (id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
		 name varchar,
		 quantity int,
		 tags text[],
		 pickup_loc geography)`)

	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `CREATE INDEX IF NOT EXISTS product_gindx ON product USING GIST (pickup_loc)`)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// InitStoreInCassandra initializes data model for seller service in the given cassandra server.
func InitStoreInCassandra(ctx context.Context, clusterHosts []string, keyspace string) error {
	cluster := gocql.NewCluster(clusterHosts...)
	cluster.Keyspace = keyspace
	session, err := cluster.CreateSession()
	if err != nil {
		return err
	}

	q := session.Query(createProductCaasandraQuery)
	q = q.WithContext(ctx)
	return q.Exec()
}
