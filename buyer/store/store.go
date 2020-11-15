// Package store has CRUD functionalities for seller service
package store

import (
	"context"
	"time"

	"github.com/gocql/gocql"
	"github.com/golang/glog"
	"github.com/ihtkas/farm/account"
	buyerpb "github.com/ihtkas/farm/buyer/v1"
)

// Storage provides a persistent storage for Seller service
type Storage struct {
	cassandra *gocql.Session
}

// OrderProduct adds a new order for the user set in context
func (s *Storage) OrderProduct(ctx context.Context, info *buyerpb.OrderInfo) error {

	uuid, err := gocql.RandomUUID()
	if err != nil {
		glog.Errorln("Error in scan after insert", err)
		return err
	}

	query := s.cassandra.Query(insertProductCassandraQuery,
		uuid,
		info.ProductId,
		info.Quantity,
		info.LocLat,
		info.LocLon,
	)

	query = query.WithContext(ctx)
	err = query.Exec()
	if err != nil {
		// TODO: Check how to retry this or rollback changes in postgres. problem in hybrid storage ie., duplicating state in two stores.
		glog.Errorln("Error in insert into cassandra:", err, insertProductCassandraQuery)
		return err
	}

	timeUUID := gocql.UUIDFromTime(time.Now())
	userID := ctx.Value(account.UserIDKey)
	query = s.cassandra.Query(insertUserOrderCassandraQuery,
		userID,
		timeUUID,
		uuid,
		info.ProductId,
		info.Quantity,
		info.LocLat,
		info.LocLon,
	)

	query = query.WithContext(ctx)
	err = query.Exec()
	if err != nil {
		glog.Errorln("Error in insert into cassandra:", err)
		return err
	}
	return nil

}

// GetOrdersList returns the Orders for specific user and the last timestamp uuid. Last timestamp uuid will be helpful for pagination.
func (s *Storage) GetOrdersList(ctx context.Context, req *buyerpb.OrdersByUserRequest) ([]*buyerpb.Order, string, error) {
	userID := ctx.Value(account.UserIDKey)
	var q *gocql.Query
	if req.LastTimestampUuid == "" {
		q = s.cassandra.Query(selectOrderByUserCassandraQuery, userID, req.Limit)
	} else {
		q = s.cassandra.Query(selectOrderByUserAfterTimeCassandraQuery, userID, req.LastTimestampUuid, req.Limit)
	}
	q = q.WithContext(ctx)
	it := q.Iter()
	defer it.Close()
	var timestampUUID string
	var orders []*buyerpb.Order
	flag := true
	for flag {
		info := &buyerpb.OrderInfo{}
		order := &buyerpb.Order{Info: info}
		flag = it.Scan(&timestampUUID, &order.Id, &info.ProductId,
			&info.Quantity, &info.LocLat, &info.LocLon)

		orders = append(orders, order)
	}
	err := it.Close()
	if err != nil {
		glog.Errorln(err)
		return nil, "", err
	}
	return orders, timestampUUID, nil
}

// Init intializes a connection pool with given postgres instance.
// TODO: POstgres intances will be scalled based on region based partitions and ht correct postgres cluster is indentified through config service in runtime
func (s *Storage) Init(ctx context.Context, cassandraClusterHosts []string, cassandraKeyspace string) error {

	cluster := gocql.NewCluster(cassandraClusterHosts...)
	cluster.Keyspace = cassandraKeyspace
	session, err := cluster.CreateSession()
	if err != nil {
		return err
	}
	s.cassandra = session
	return nil
}

// id uuid, minquantity int, price_per_quantity int, description varchar,
// seller_id uuid, expiry interval,

// CREATE TABLE IF NOT EXISTS product (id uuid DEFAULT uuid_generate_v4() PRIMARY KEY, name varchar, quantity int, tags text[], pickup_loc geography);

// select id, ST_Distance(pickup_loc, ref_geoloc) AS distance
// from product
// CROSS JOIN (SELECT ST_MakePoint(11.530012, 78.108304)::geography AS ref_geoloc) AS r
// WHERE ST_DWithin(pickup_loc, ref_geoloc, 1000)
// ORDER BY ST_Distance(pickup_loc, ref_geoloc);
