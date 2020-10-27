// Package store has CRUD functionalities for seller service
package store

import (
	"context"
	"fmt"

	"github.com/gocql/gocql"
	"github.com/golang/glog"
	sellerpb "github.com/ihtkas/farm/seller/v1"
	pgpoolv4 "github.com/jackc/pgx/v4/pgxpool"
)

// Storage provides a persistent storage for Seller service
type Storage struct {
	pg        *pgpoolv4.Pool
	cassandra *gocql.Session
}

// AddProduct adds a new product
func (s *Storage) AddProduct(ctx context.Context, p *sellerpb.ProductInput) error {

	rows, err := s.pg.Query(ctx,
		`INSERT INTO product (name, quantity, tags, pickup_loc) values ($1, $2, $3, $4) returning (id)
	`, p.Name, p.Quantity, p.Tags, fmt.Sprintf("Point(%v %v)", p.PickupLocLat, p.PickupLocLon))
	if err != nil {
		return err
	}
	var id string
	err = rows.Scan(&id)
	if err != nil {
		return err
	}
	query := s.cassandra.Query(`
	INSERT INTO product (
		id,
		name,
		expiry,
		quantity,
		minquantity,
		priceperquantity,
		description,
		tags) values (?, ?, ?, ?, ?, ?, ?, ?)	
	`,
		id,
		p.Name,
		p.Expiry,
		p.Quantity,
		p.MinQuantity,
		p.PricePerQuantity,
		p.Description,
		p.Tags,
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
func (s *Storage) GetNearbyProducts(ctx context.Context, lat, lon float64, radius, offset, limit int) ([]*sellerpb.ProductResponse, error) {
	products, err := s.getNearbyProductsPG(ctx, lat, lon, radius, offset, limit)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (s *Storage) getNearbyProductsPG(ctx context.Context, lat, lon float64, radius, offset, limit int) ([]*sellerpb.ProductResponse, error) {

	conn, err := s.pg.Acquire(ctx)
	if err != nil {
		return nil, err
	}

	// uses postgis extension  query to sort nearby location .
	rows, err := conn.Query(ctx, `
	SELECT id, ST_Distance(pickup_loc, ref_geoloc) AS distance
	FROM product CROSS JOIN (SELECT ST_MakePoint($1, $2)::geography AS ref_geoloc) AS r
	WHERE ST_DWithin(pickup_loc, ref_geoloc, $3)
	ORDER BY ST_Distance(pickup_loc, ref_geoloc) LIMIT $4, OFFSET $5`,
		lat, lon, radius, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	products := make([]*sellerpb.ProductResponse, 0)
	// get full info about the product form cassandra
	for rows.Next() {
		p := &sellerpb.ProductResponse{}
		err := rows.Scan(&p.Id, &p.Distance)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func (s *Storage) fetchProductDetails(ctx context.Context, products []*sellerpb.ProductResponse) error {
	for _, p := range products {
		q := s.cassandra.Query("SELECT name, expiry, quantity, min_quantity, price_per_quantity, description, tags, pickup_loc_lat, pickup_loc_lon from product where id=?", p.Id)
		q = q.WithContext(ctx)
		err := q.Exec()
		if err != nil {
			return err
		}

		err = q.Scan(&p.Name, &p.Expiry, &p.Quantity, &p.MinQuantity, &p.PricePerQuantity, &p.Description, &p.Tags, &p.PickupLocLat, &p.PickupLocLon)
		if err != nil {
			return err
		}
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
