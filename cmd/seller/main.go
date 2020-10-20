package main

import (
	"github.com/golang/glog"
	"github.com/ihtkas/farm/seller"
	"github.com/ihtkas/farm/seller/store"
)

func main() {
	store := &store.Cassandra{}
	err := store.Init([]string{"127.0.0.1"}, "farm")
	if err != nil {
		glog.Errorln(err)
		return
	}
	var broker seller.MessageProducer
	m := &seller.Manager{}
	err = m.Start(store, broker)
	if err != nil {
		glog.Errorln(err)
		return
	}
	glog.Errorln("Done....")
}
