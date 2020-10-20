// Package store has CRUD functionalities for seller service
package store

import (
	"errors"

	"github.com/gocql/gocql"
	accountpb "github.com/ihtkas/farm/account/v1"
)

// Cassandra provides a persistent storage for Seller service
type Cassandra struct {
	session *gocql.Session
}

// AddUser adds a new user
func (s *Cassandra) AddUser(p *accountpb.User) error {
	id, err := gocql.RandomUUID()
	if err != nil {
		return err
	}
	return s.session.Query(`
	INSERT INTO user (id, name) values (?, ?)	
	`, id, p.Name,
	).Exec()
}

// ValidateUser checks the existence of user in database
func (s *Cassandra) ValidateUser(p *accountpb.ValidateUserRequest) error {
	return errors.New("Not implemented")
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
	CREATE TABLE IF NOT EXISTS user (
		id UUID PRIMARY KEY, 
		name varchar
	)
	`).Exec()

}
