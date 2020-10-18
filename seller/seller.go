package seller

import (
	"fmt"
	"io"
	"net/http"
	"net/http/pprof"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/ihtkas/farm/seller/store"
	"github.com/ihtkas/farm/utils"

	"github.com/golang/glog"
	sellerpb "github.com/ihtkas/farm/seller/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

// Manager implements the http.Handler interface and manages all APIs for account management
type Manager struct {
	addr                  string
	cassandraClusterHosts []string
	cassandraKeyspace     string
	store                 Storage

	minExpiryDur time.Duration
}

// Storage has functions required to store, read and manipulate Seller information
type Storage interface {
	Init(clusterHosts []string, keyspace string) error
	AddProduct(p *sellerpb.Product) error
}

func (m *Manager) initDefaultConf() {
	m.addr = ":8082"
	m.cassandraClusterHosts = []string{"127.0.0.1"}
	m.cassandraKeyspace = "farm"
	m.store = &store.Cassandra{}
	m.minExpiryDur = 12 * time.Hour

}

// Start first initializes with default configuration and overrides with input options. Then starts a http server.
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
	// 	// fmt.Println(reflect.TypeOf(d.Values[0]), *(d.Values[0].(*string)))
	// 	fmt.Println(s1, s2)
	// }
	s := &http.Server{Addr: m.addr, Handler: m}
	return s.ListenAndServe()

}

// ServeHTTP handles all http APIs for Account management
func (m *Manager) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	glog.Errorln(r.URL.Path)
	switch r.URL.Path {
	case "/seller/product/add":
		m.addProductReq(w, r)
	default:
		if strings.HasPrefix(r.URL.Path, "/debug") {
			arr := strings.Split(r.URL.Path, "/")
			if len(arr) >= 4 {
				glog.Errorln("pprof", arr[3])
				pprof.Handler(arr[3]).ServeHTTP(w, r)
				return
			}
		}
		http.Error(w, "Invalid path", http.StatusBadRequest)
	}
}

func (m *Manager) addProductReq(w http.ResponseWriter, r *http.Request) {

	body := make([]byte, 1<<10)
	n, err := r.Body.Read(body)

	glog.Errorln(string(body[:n]), n, err)
	if err != nil && err != io.EOF {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	product := &sellerpb.Product{}

	err = protojson.Unmarshal(body[:n], product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println(product)
	minExpiry := time.Now().Add(m.minExpiryDur)
	validation.ValidateStruct(&product,
		validation.Field(&product.Name, validation.Required, validation.Length(5, 50)),
		validation.Field(&product.Expiry, validation.Required, utils.TimeRange(minExpiry, time.Time{})),
		validation.Field(&product.Quantity, validation.Required),
		validation.Field(&product.PricePerQuantity, validation.Required),
		validation.Field(&product.MinQuantity, validation.Required),
	)

	err = m.store.AddProduct(product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
