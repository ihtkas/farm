package buyer

import (
	"context"
	"io"
	"net/http"
	"net/http/pprof"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/ihtkas/farm/account"
	"github.com/ihtkas/farm/buyer/store"
	buyerpb "github.com/ihtkas/farm/buyer/v1"
	"github.com/ihtkas/farm/utils"

	"github.com/golang/glog"
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

// Storage has functions required to store, read and manipulate buyer information
type Storage interface {
	OrderProduct(ctx context.Context, p *buyerpb.OrderInfo) error
	GetOrdersList(ctx context.Context, req *buyerpb.OrdersByUserRequest) ([]*buyerpb.Order, string, error)
}

// MessageProducer has functions to publish new products to the matching system
// TODO: explore Kafka connect for cassandra instead of manual publish
type MessageProducer interface {
	PublishNewOrder(p *buyerpb.OrderInfo) error
}

func (m *Manager) initDefaultConf() {
	m.addr = ":8081"
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
	// TODO: get the userID from cookie and set it in context
	ctx := context.WithValue(r.Context(), account.UserIDKey, "1105ce66-95da-4e5b-9af2-25976c8f4f5d")
	r = r.WithContext(ctx)
	// ---------------

	switch r.URL.Path {
	case "/buyer/order/add":
		m.addOrdeReq(w, r)
	case "/buyer/order/list":
		m.getOrdersList(w, r)
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

func (m *Manager) addOrdeReq(w http.ResponseWriter, r *http.Request) {

	body := make([]byte, 1<<10)
	n, err := r.Body.Read(body)

	glog.Errorln(string(body[:n]), n, err)
	if err != nil && err != io.EOF {
		glog.Errorln(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	order := &buyerpb.OrderInfo{}

	err = protojson.Unmarshal(body[:n], order)
	if err != nil {
		glog.Errorln(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	validation.ValidateStruct(&order,
		validation.Field(&order.ProductId, validation.Required, validation.Length(5, 50)),
		validation.Field(&order.Quantity, validation.Required),
		validation.Field(&order.LocLat, validation.Required),
		validation.Field(&order.LocLon, validation.Required),
	)

	err = m.store.OrderProduct(r.Context(), order)
	if err != nil {
		glog.Errorln(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// err = m.msgBroker.PublishNewOrder(order)
	// if err != nil {
	// 	// TODO: handle this failure case. The order has to be injected again into the matcher module
	// 	// May try Kafka connect to pull directly from store (cassandra)
	// 	// Add alerts for such failures.
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }
}

func (m *Manager) getOrdersList(w http.ResponseWriter, r *http.Request) {

	body := make([]byte, 1<<10)
	n, err := r.Body.Read(body)

	if err != nil && err != io.EOF {
		glog.Errorln(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	req := &buyerpb.OrdersByUserRequest{}

	err = protojson.Unmarshal(body[:n], req)
	if err != nil {
		glog.Errorln(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	validation.ValidateStruct(&req,
		validation.Field(&req.Limit, utils.Uint32Range(limitMin, limitMax, true, true)),
	)

	orders, uuid, err := m.store.GetOrdersList(r.Context(), req)
	if err != nil {
		glog.Errorln(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var arr []string
	for _, p := range orders {
		b, err := protojson.Marshal(p)
		if err != nil {
			glog.Errorln(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		arr = append(arr, string(b))
	}
	_, err = w.Write([]byte(`{"orders": [` + strings.Join(arr, `, `) + `], "last_timeUUID": "` + uuid + `"}`))
	if err != nil {
		glog.Errorln(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// err = m.msgBroker.PublishNeworder(order)
	// if err != nil {
	// 	// TODO: handle this failure case. The order has to be injected again into the matcher module
	// 	// May try Kafka connect to pull directly from store (cassandra)
	// 	// Add alerts for such failures.
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }
}
