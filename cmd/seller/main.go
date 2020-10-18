package main

import (
	"github.com/golang/glog"
	"github.com/ihtkas/farm/seller"
)

func main() {
	m := seller.Manager{}
	err := m.Start()
	if err != nil {
		glog.Errorln(err)
		return
	}
	glog.Errorln("Done....")
}
