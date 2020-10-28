package seller

import (
	"context"
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

const (
	radiusMin = 1     // 1 metre
	radiusMax = 20000 // 20 km
	offsetMin = 0
	offsetMax = 100
	limitMin  = 1
	limitMax  = 50
)

// Manager implements the http.Handler interface and manages all APIs for account management
type Manager struct {
	addr         string
	store        Storage
	msgBroker    MessageProducer
	minExpiryDur time.Duration
}

// Storage has functions required to store, read and manipulate Seller information
type Storage interface {
	AddProduct(ctx context.Context, p *sellerpb.ProductInput) error
	GetNearbyProducts(ctx context.Context, loc *sellerpb.ProductLocationSearchRequest) ([]*sellerpb.ProductResponse, error)
}

// MessageProducer has functions to publish new products to the matching system
// TODO: explore Kafka connect for cassandra instead of manual publish
type MessageProducer interface {
	PublishNewProduct(p *sellerpb.ProductInput) error
}

func (m *Manager) initDefaultConf() {
	m.addr = ":8082"
	m.store = &store.Storage{}
	m.minExpiryDur = 12 * time.Hour
}

// Start first initializes with default configuration and overrides with input options. Then starts a http server.
func (m *Manager) Start(store Storage, broker MessageProducer, opts ...Option) error {
	m.initDefaultConf()

	for _, opt := range opts {
		err := opt.set(m)
		if err != nil {
			return err
		}
	}
	m.store = store
	m.msgBroker = broker
	// iter := session.Query("SELECT cluster_name, listen_address FROM system.local;").Iter()
	// var s1, s2 string
	// exist := iter.Scan(&s1, &s2)
	// if exist {
	// 	// fmt.Println(reflect.TypeOf(d.Values[0]), *(d.Values[0].(*string)))
	// 	fmt.Println(s1, s2)
	// }
	s := &http.Server{Addr: m.addr, Handler: m}
	glog.V(0).Infoln("Starting server in", m.addr)
	return s.ListenAndServe()

}

// ServeHTTP handles all http APIs for Account management
func (m *Manager) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/seller/product/add":
		m.addProductReq(w, r)
	case "/seller/product/get_nearby_product":
		m.getNearByProduct(w, r)
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
		glog.Errorln(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	product := &sellerpb.ProductInput{}

	err = protojson.Unmarshal(body[:n], product)
	if err != nil {
		glog.Errorln(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	minExpiry := time.Now().Add(m.minExpiryDur)
	validation.ValidateStruct(&product,
		validation.Field(&product.Name, validation.Required, validation.Length(5, 50)),
		validation.Field(&product.Expiry, validation.Required, utils.TimeRange(minExpiry, time.Time{})),
		validation.Field(&product.Quantity, validation.Required),
		validation.Field(&product.PricePerQuantity, validation.Required),
		validation.Field(&product.MinQuantity, validation.Required),
		validation.Field(&product.PickupLocLat, validation.Required),
		validation.Field(&product.PickupLocLon, validation.Required),
	)

	err = m.store.AddProduct(r.Context(), product)
	if err != nil {
		glog.Errorln(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// err = m.msgBroker.PublishNewProduct(product)
	// if err != nil {
	// 	// TODO: handle this failure case. The product has to be injected again into the matcher module
	// 	// May try Kafka connect to pull directly from store (cassandra)
	// 	// Add alerts for such failures.
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }
}

func (m *Manager) getNearByProduct(w http.ResponseWriter, r *http.Request) {

	body := make([]byte, 1<<10)
	n, err := r.Body.Read(body)

	if err != nil && err != io.EOF {
		glog.Errorln(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	loc := &sellerpb.ProductLocationSearchRequest{}

	err = protojson.Unmarshal(body[:n], loc)
	if err != nil {
		glog.Errorln(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	validation.ValidateStruct(&loc,
		validation.Field(&loc.PickupLocLat, utils.Float64Range(-90, 90, true, true)),
		validation.Field(&loc.PickupLocLon, utils.Float64Range(-180, 180, true, true)),
		validation.Field(&loc.Radius, utils.Uint32Range(radiusMin, radiusMax, true, true)),
		validation.Field(&loc.Offset, utils.Uint32Range(offsetMin, offsetMax, true, true)),
		validation.Field(&loc.Limit, utils.Uint32Range(limitMin, limitMax, true, true)),
	)

	products, err := m.store.GetNearbyProducts(r.Context(), loc)
	if err != nil {
		glog.Errorln(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var arr []string
	for _, p := range products {
		b, err := protojson.Marshal(p)
		if err != nil {
			glog.Errorln(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		arr = append(arr, string(b))
	}
	_, err = w.Write([]byte("[" + strings.Join(arr, ", ") + "]"))
	if err != nil {
		glog.Errorln(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// err = m.msgBroker.PublishNewProduct(product)
	// if err != nil {
	// 	// TODO: handle this failure case. The product has to be injected again into the matcher module
	// 	// May try Kafka connect to pull directly from store (cassandra)
	// 	// Add alerts for such failures.
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }
}
