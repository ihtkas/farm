package seller

// Option overrides default configuration of the Account Manager
type Option interface {
	// contains filtered or unexported methods
	set(m *Manager) error
}

// HTTPAddr retuns an option to set http server address
func HTTPAddr(addr string) Option {
	return &httpAddrConfOpt{addr: addr}
}

type httpAddrConfOpt struct {
	addr string
}

func (opt *httpAddrConfOpt) set(m *Manager) error {
	m.addr = opt.addr
	return nil
}

// CassandraClusterHosts retuns an option to set casscandra cluster hosts
func CassandraClusterHosts(hosts []string) Option {
	return &cassandraClusterHostsConfOpt{hosts: hosts}
}

type cassandraClusterHostsConfOpt struct {
	hosts []string
}

func (opt *cassandraClusterHostsConfOpt) set(m *Manager) error {
	m.cassandraClusterHosts = opt.hosts
	return nil
}

// CassandraKeySpace retuns an option to set casscandra key space
func CassandraKeySpace(keyspace string) Option {
	return &cassandraKeySpaceConfOpt{keyspace: keyspace}
}

type cassandraKeySpaceConfOpt struct {
	keyspace string
}

func (opt *cassandraKeySpaceConfOpt) set(m *Manager) error {
	m.cassandraKeyspace = opt.keyspace
	return nil
}
