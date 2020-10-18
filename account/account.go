package account

import (
	"io"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/golang/glog"
	accountpb "github.com/ihtkas/farm/account/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

// Manager implements the http.Handler interface and manages all APIs for account management
type Manager struct {
	addr                  string
	cassandraClusterHosts []string
	cassandraKeyspace     string
	store                 Storage
}

// Storage has functions required to store, read and manipulate Seller information
type Storage interface {
	Init(clusterHosts []string, keyspace string) error
	AddUser(p *accountpb.User) error
	ValidateUser(p *accountpb.ValidateUserRequest) error
}

func (m *Manager) initDefaultConf() {
	m.addr = ":8081"
	m.cassandraClusterHosts = []string{"127.0.0.1"}
	m.cassandraKeyspace = "farm"
}

// Start starts a Account Manager service
func (m *Manager) Start(opts ...Option) error {
	m.initDefaultConf()

	for _, opt := range opts {
		err := opt.set(m)
		if err != nil {
			return err
		}
	}
	err := m.store.Init(m.cassandraClusterHosts, m.cassandraKeyspace)
	if err != nil {
		return err
	}

	// iter := session.Query("SELECT cluster_name, listen_address FROM system.local;").Iter()
	// var s1, s2 string
	// exist := iter.Scan(&s1, &s2)
	// if exist {
	// }
	s := &http.Server{Addr: m.addr, Handler: m}
	return s.ListenAndServe()

}

// ServeHTTP handles all http APIs for Account management
func (m *Manager) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	switch r.URL.Path {
	case "account/user/add":
		m.addUserReq(w, r)
	default:
		http.Error(w, "Invalid path", http.StatusBadRequest)
	}
}

func (m *Manager) addUserReq(w http.ResponseWriter, r *http.Request) {
	body := make([]byte, 1<<10)
	n, err := r.Body.Read(body)

	glog.Errorln(string(body[:n]), n, err)
	if err != nil && err != io.EOF {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user := &accountpb.User{}

	err = protojson.Unmarshal(body[:n], user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	validation.ValidateStruct(&user,
		validation.Field(&user.Name, validation.Required, validation.Length(5, 50)),
	)

	err = m.store.AddUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
