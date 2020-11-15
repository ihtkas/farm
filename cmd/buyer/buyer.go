package main

import (
	"context"
	"flag"

	"github.com/golang/glog"
	"github.com/ihtkas/farm/buyer"
	"github.com/ihtkas/farm/buyer/store"
	"github.com/ihtkas/farm/utils"
)

var cassandraClusterHosts utils.ArrayFlags
var cassandraKeyspace string

func main() {
	flag.Var(&cassandraClusterHosts, "cassandra_cluster_hosts", "List of hosts for cassandra hosts")
	flag.StringVar(&cassandraKeyspace, "cassandra_cluster_keyspace", "farm_buyer", "Keyaspace in cassandra for buyer service")
	flag.Parse()
	// set default configuration if non provided for development convenience.
	if len(cassandraClusterHosts) == 0 {
		err := cassandraClusterHosts.Set("127.0.0.1:9042")
		if err != nil {
			glog.Errorln(err)
			return
		}
	}

	store := &store.Storage{}
	err := store.Init(context.Background(), cassandraClusterHosts, cassandraKeyspace)
	if err != nil {
		glog.Errorln(err)
		return
	}
	var broker buyer.MessageProducer
	m := &buyer.Manager{}
	err = m.Start(store, broker)
	if err != nil {
		glog.Errorln(err)
		return
	}
	glog.Errorln("Done....")
}
