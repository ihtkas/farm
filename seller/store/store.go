// Package store has CRUD functionalities for seller service
package store

import (
	"github.com/gocql/gocql"
	sellerpb "github.com/ihtkas/farm/seller/v1"
)

// Cassandra provides a persistent storage for Seller service
type Cassandra struct {
	session *gocql.Session
}

// AddProduct adds a new product
func (s *Cassandra) AddProduct(p *sellerpb.Product) error {
	return s.session.Query(`
	INSERT INTO product (
		name,
		expiry,
		quantity,
		minquantity,
		priceperquantity,
		description,
		tags) values (?, ?)	
	`,
		p.Name,
		p.Expiry,
		p.Quantity,
		p.MinQuantity,
		p.PricePerQuantity,
		p.Description,
		p.Tags,
	).Exec()
}

// Init intializes a cassandra storage with given configuration.
func (s *Cassandra) Init(clusterHosts []string, keyspace string) error {
	cluster := gocql.NewCluster(clusterHosts...)
	cluster.Keyspace = keyspace
	session, err := cluster.CreateSession()
	if err != nil {
		return err
	}
	s.session = session
	return session.Query(`
	CREATE TABLE IF NOT EXISTS product (
		id UUID PRIMARY KEY, 
		name varchar,
		expiry duration,
		quantity int,
		minquantity int,
		price_per_quantity int,
		description varchar,
		tags list<varchar>,
	)
	`).Exec()

}
