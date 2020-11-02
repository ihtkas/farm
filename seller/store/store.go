// Package store has CRUD functionalities for seller service
package store

import (
	"context"
	"fmt"
	"time"

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
func (s *Storage) AddProduct(ctx context.Context, info *sellerpb.ProductInfo) error {

	rows, err := s.pg.Query(ctx,
		insertProductPGQuery, info.Name, info.Quantity, info.Tags, fmt.Sprintf("Point(%v %v)", info.PickupLocLat, info.PickupLocLon))
	if err != nil {
		glog.Errorln("Error in insert into postgres", err, insertProductPGQuery)
		return err
	}
	var id string
	rows.Next()
	err = rows.Scan(&id)
	if err != nil {
		glog.Errorln("Error in scan after insert", err)
		return err
	}
	infoBlob, err := proto.Marshal(info)
	if err != nil {
		glog.Errorln("Error in proto marshal", err)
		return err
	}
	query := s.cassandra.Query(insertProductCassandraQuery,
		id,
		infoBlob,
	)

	query = query.WithContext(ctx)
	err = query.Exec()
	if err != nil {
		// TODO: Check how to retry this or rollback changes in postgres. problem in hybrid storage ie., duplicating state in two stores.
		glog.Errorln("Error in insert into cassandra:", err)
		return err
	}
	product := &sellerpb.Product{Id: id, Info: info}

	productBlob, err := proto.Marshal(product)
	if err != nil {
		glog.Errorln("Error in proto marshal", err)
		return err
	}
	glog.Errorln(product, productBlob)
	time_uuid := gocql.UUIDFromTime(time.Now())
	userID := ctx.Value("UserID")
	query = s.cassandra.Query(insertUserProductCassandraQuery,
		userID,
		time_uuid,
		productBlob,
	)

	query = query.WithContext(ctx)
	err = query.Exec()
	if err != nil {
		// TODO: Check how to retry this or rollback changes in postgres. problem in hybrid storage ie., duplicating state in two stores.
		glog.Errorln("Error in insert into cassandra:", err)
		return err
	}
	return nil

}

// GetNearbyProducts finds the products sent form pickup locations
// within given radius ordered by shortest distances first.
func (s *Storage) GetNearbyProducts(ctx context.Context, loc *sellerpb.ProductLocationSearchRequest) ([]*sellerpb.ProductLocationResponse, error) {
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
func (s *Storage) getNearbyProductsPostgres(ctx context.Context, loc *sellerpb.ProductLocationSearchRequest) ([]*sellerpb.ProductLocationResponse, error) {

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
	products := make([]*sellerpb.ProductLocationResponse, 0)
	// get full info about the product form cassandra
	for rows.Next() {
		p := &sellerpb.ProductLocationResponse{}
		err := rows.Scan(&p.Id, &p.Distance)
		if err != nil {
			glog.Errorln(err)
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func (s *Storage) fetchProductDetails(ctx context.Context, products []*sellerpb.ProductLocationResponse) error {
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

// GetProductsList returns the products for specific user and the last timestamp uuid. Last timestamp uuid will be helpful for pagination.
func (s *Storage) GetProductsList(ctx context.Context, req *sellerpb.ProductsByUserRequest) ([]*sellerpb.Product, string, error) {
	userID := ctx.Value("UserID")
	var q *gocql.Query
	if req.LastTimestampUuid == "" {
		q = s.cassandra.Query(selectProductByUserCassandraQuery, userID, req.Limit)
	} else {
		q = s.cassandra.Query(selectProductByUserCassandraQuery, userID, req.LastTimestampUuid, req.Limit)
	}
	q = q.WithContext(ctx)
	it := q.Iter()
	defer it.Close()
	var productBlob []byte
	var timestampUUID string
	var products []*sellerpb.Product
	for it.Scan(&productBlob, &timestampUUID) {
		product := &sellerpb.Product{}
		err := proto.Unmarshal(productBlob, product)
		if err != nil {
			glog.Errorln(err, q, timestampUUID, productBlob)
			return nil, "", err
		}
		products = append(products, product)
	}
	err := it.Close()
	if err != nil {
		glog.Errorln(err)
		return nil, "", err
	}
	return products, timestampUUID, nil
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
