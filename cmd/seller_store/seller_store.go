package seller_store

import (
	"context"
	"flag"

	"github.com/golang/glog"

	"github.com/ihtkas/farm/utils"

	"github.com/ihtkas/farm/seller/store"
)

var cassandraClusterHosts utils.ArrayFlags
var cassandraKeyspace string
var pgUsername, pgPassword, pgHost, pgPort, pgDbname string

func main() {
	//  clusterHosts []string, keyspace string
	flag.Var(&cassandraClusterHosts, "cassandra_cluster_hosts", "List of hosts for cassandra hosts")
	flag.StringVar(&cassandraKeyspace, "cassandra_cluster_keyspace", "farm_seller", "Keyaspace in cassandra for seller service")
	flag.StringVar(&pgUsername, "postgres_user_name", "postgres", "User to connect postgres server")
	flag.StringVar(&pgPassword, "postgres_password", "", "Password to connect postrgres server")
	flag.StringVar(&pgHost, "postgres_host", "127.0.0.1", "Host address of postgres server")
	flag.StringVar(&pgPort, "postgres_port", "54321", "Port in the host to connect postgres server")
	flag.StringVar(&pgDbname, "postgres_dbname", "farm_seller", "Database name for seller service in postgres")

	ctx := context.Background()
	err := store.InitStoreInCassandra(ctx, cassandraClusterHosts, cassandraKeyspace)
	if err != nil {
		glog.Errorln(err)
		return
	}
	flag.StringVar(&cassandraKeyspace, "Cassandra Cluster Hosts", "farm_seller", "List of hosts for cassandra hosts")

	err = store.InitStoreInPostgres(ctx, pgUsername, pgPassword, pgHost, pgPort, pgDbname)
	if err != nil {
		glog.Errorln(err)
		return
	}

}
