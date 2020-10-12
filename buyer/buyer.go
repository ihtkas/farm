package buyer

import (
	"net/http"

	"github.com/gocql/gocql"
	"github.com/ihtkas/farm/utils"
)

// Manager implements the http.Handler interface and manages all APIs for account management
type Manager struct {
	addr                  string
	cassandraClusterHosts []string
	cassandraKeyspace     string
	session               *gocql.Session
}

func (m *Manager) initDefaultConf() {
	m.addr = ":8083"
	m.cassandraClusterHosts = []string{"127.0.0.1"}
	m.cassandraKeyspace = "farm"
}

// Starts an Manager service
func (m *Manager) Start(opts ...Option) error {
	m.initDefaultConf()

	for _, opt := range opts {
		err := opt.set(m)
		if err != nil {
			return err
		}
	}
	err := m.initStorage()
	if err != nil {
		return err
	}
	// iter := session.Query("SELECT cluster_name, listen_address FROM system.local;").Iter()
	// var s1, s2 string
	// exist := iter.Scan(&s1, &s2)
	// if exist {
	// 	// fmt.Println(reflect.TypeOf(d.Values[0]), *(d.Values[0].(*string)))
	// 	fmt.Println(s1, s2)
	// }
	s := &http.Server{Addr: m.addr, Handler: m}
	return s.ListenAndServe()

}

func (m *Manager) initStorage() error {
	cluster := gocql.NewCluster(m.cassandraClusterHosts...)
	cluster.Keyspace = m.cassandraKeyspace
	session, err := cluster.CreateSession()
	if err != nil {
		return err
	}
	m.session = session
	return session.Query(`
	CREATE TABLE IF NOT EXISTS user (
		id int PRIMARY KEY, 
		name varchar
	)
	`).Exec()

}

// ServeHTTP handles all http APIs for Account management
func (m *Manager) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	switch r.URL.Path {
	case "/user/add":
		m.addUserReq(w, r)
	default:
		http.Error(w, "Invalid path", http.StatusBadRequest)
	}
}

func (m *Manager) addUserReq(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	username, err := utils.GetStringParam(r.Form, "username")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := utils.GetIntegerParam(r.Form, "id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = m.addUser(username, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (m *Manager) addUser(name string, id int64) error {
	return m.session.Query(`
	INSERT INTO user (id, name) values (?, ?)	
	`, id, name,
	).Exec()
}
