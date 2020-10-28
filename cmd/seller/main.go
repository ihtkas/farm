package main

import (
	"context"
	"flag"

	"github.com/golang/glog"
	"github.com/ihtkas/farm/seller"
	"github.com/ihtkas/farm/seller/store"
	"github.com/ihtkas/farm/utils"
)

var cassandraClusterHosts utils.ArrayFlags
var cassandraKeyspace string
var pgUsername, pgPassword, pgHost, pgPort, pgDbname string

func main() {
	flag.Var(&cassandraClusterHosts, "cassandra_cluster_hosts", "List of hosts for cassandra hosts")
	flag.StringVar(&cassandraKeyspace, "cassandra_cluster_keyspace", "farm_seller", "Keyaspace in cassandra for seller service")
	flag.StringVar(&pgUsername, "postgres_user_name", "postgres", "User to connect postgres server")
	flag.StringVar(&pgPassword, "postgres_password", "", "Password to connect postrgres server")
	flag.StringVar(&pgHost, "postgres_host", "127.0.0.1", "Host address of postgres server")
	flag.StringVar(&pgPort, "postgres_port", "5432", "Port in the host to connect postgres server")
	flag.StringVar(&pgDbname, "postgres_dbname", "farm_seller", "Database name for seller service in postgres")

	// set default configuration if non provided for development convenience.
	if len(cassandraClusterHosts) == 0 {
		err := cassandraClusterHosts.Set("127.0.0.1:9042")
		if err != nil {
			glog.Errorln(err)
			return
		}
	}

	store := &store.Storage{}
	err := store.Init(context.Background(), pgUsername, pgPassword, pgHost, pgPort, pgDbname, cassandraClusterHosts, cassandraKeyspace)
	if err != nil {
		glog.Errorln(err)
		return
	}
	var broker seller.MessageProducer
	m := &seller.Manager{}
	err = m.Start(store, broker)
	if err != nil {
		glog.Errorln(err)
		return
	}
	glog.Errorln("Done....")
}
