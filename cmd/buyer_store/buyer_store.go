package main

import (
	"context"
	"flag"

	"github.com/golang/glog"

	"github.com/ihtkas/farm/utils"

	"github.com/ihtkas/farm/buyer/store"
)

var cassandraClusterHosts utils.ArrayFlags
var cassandraKeyspace string

func main() {
	// Prerequisites for development testing:
	// Cassandra:
	// CREATE KEYSPACE farm_buyer WITH replication = {'class':'SimpleStrategy', 'replication_factor' : 1};
	// Postgres:
	// CREATE DATABSE farm_buyer;

	flag.Var(&cassandraClusterHosts, "cassandra_cluster_hosts", "List of hosts for cassandra hosts")
	flag.StringVar(&cassandraKeyspace, "cassandra_cluster_keyspace", "farm_buyer", "Keyspace in cassandra for buyer service")
	flag.Parse()
	// set default configuration if non provided for development convenience.
	if len(cassandraClusterHosts) == 0 {
		err := cassandraClusterHosts.Set("127.0.0.1:9042")
		if err != nil {
			glog.Errorln(err)
			return
		}
	}
	ctx := context.Background()
	err := store.InitStoreInCassandra(ctx, cassandraClusterHosts, cassandraKeyspace)
	if err != nil {
		glog.Errorln(err)
		return
	}

}
