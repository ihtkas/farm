// Package store has CRUD functionalities for seller service
package store

import (
	"context"
	"fmt"

	"github.com/gocql/gocql"
	"github.com/golang/glog"
	sellerpb "github.com/ihtkas/farm/seller/v1"
	pgpoolv4 "github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/protobuf/proto"
)

// Storage provides a persistent storage for Seller service
type Storage struct {
	pg        *pgpoolv4.Pool
	cassandra *gocql.Session
}

// AddProduct adds a new product
func (s *Storage) AddProduct(ctx context.Context, product *sellerpb.ProductInfo) error {

	rows, err := s.pg.Query(ctx,
		insertProductPGQuery, product.Name, product.Quantity, product.Tags, fmt.Sprintf("Point(%v %v)", product.PickupLocLat, product.PickupLocLon))
	if err != nil {
		glog.Errorln("Error in insert into postgres", err)
		return err
	}
	var id string
	rows.Next()
	err = rows.Scan(&id)
	if err != nil {
		glog.Errorln("Error in scan after insert", err)
		return err
	}
	productBlob, err := proto.Marshal(product)
	if err != nil {
		glog.Errorln("Error in proto marshal", err)
		return err
	}
	query := s.cassandra.Query(insertProductCassandraQuery,
		id,
		productBlob,
	)
	query = query.WithContext(ctx)
	err = query.Exec()
	if err != nil {
		// TODO: Check how to retry this or rollback changes in postgres. problem in hybrid storage ie., duplicating state in two stores.
		glog.Errorln("Error in insert into cassandra:", err)
		return err
	}
	return err

}

// GetNearbyProducts finds the products sent form pickup locations
// within given radius ordered by shortest distances first.
func (s *Storage) GetNearbyProducts(ctx context.Context, loc *sellerpb.ProductLocationSearchRequest) ([]*sellerpb.ProductResponse, error) {
	products, err := s.getNearbyProductsPostgres(ctx, loc)
	if err != nil {
		return nil, err
	}
	err = s.fetchProductDetails(ctx, products)
	if err != nil {
		return nil, err
	}
	return products, nil
}
func (s *Storage) getNearbyProductsPostgres(ctx context.Context, loc *sellerpb.ProductLocationSearchRequest) ([]*sellerpb.ProductResponse, error) {

	conn, err := s.pg.Acquire(ctx)
	if err != nil {
		return nil, err
	}

	// uses postgis extension  query to sort nearby location .

	rows, err := conn.Query(ctx, nearByProductsPGQuery,
		loc.PickupLocLat, loc.PickupLocLon, loc.Radius, loc.Limit, loc.Offset)
	if err != nil {
		glog.Errorln(err)
		return nil, err
	}
	defer rows.Close()
	products := make([]*sellerpb.ProductResponse, 0)
	// get full info about the product form cassandra
	for rows.Next() {
		p := &sellerpb.ProductResponse{}
		err := rows.Scan(&p.Id, &p.Distance)
		if err != nil {
			glog.Errorln(err)
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func (s *Storage) fetchProductDetails(ctx context.Context, products []*sellerpb.ProductResponse) error {
	for _, p := range products {
		q := s.cassandra.Query(selectProductCassandraQuery, p.Id)
		q = q.WithContext(ctx)
		err := q.Exec()
		if err != nil {
			glog.Errorln(err, q.String(), p)
			return err
		}
		var productBlob []byte
		err = q.Scan(&productBlob)
		if err != nil {
			glog.Errorln(err, p)
			return err
		}
		info := &sellerpb.ProductInfo{}
		err = proto.Unmarshal(productBlob, info)
		if err != nil {
			glog.Errorln(err, p)
			return err
		}
		p.Info = info
	}
	return nil
}

// Init intializes a connection pool with given postgres instance.
// TODO: POstgres intances will be scalled based on region based partitions and ht correct postgres cluster is indentified through config service in runtime
func (s *Storage) Init(ctx context.Context, pgUsername, pgPassword, pgHost, pgPort, pgDbname string, cassandraClusterHosts []string, cassandraKeyspace string) error {

	config := "user=" + pgUsername +
		" password=" + pgPassword +
		" host=" + pgHost +
		" port=" + pgPort +
		" dbname=" + pgDbname

	pool, err := pgpoolv4.Connect(ctx, config)
	if err != nil {
		return err
	}
	s.pg = pool
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
